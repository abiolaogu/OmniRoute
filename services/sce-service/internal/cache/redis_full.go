// Package cache provides caching implementations for the pricing engine
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/omniroute/pricing-engine/internal/domain"
)

// RedisCache implements the PriceCache interface using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	if prefix == "" {
		prefix = "pricing:"
	}
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a cached price response
func (c *RedisCache) Get(ctx context.Context, key string) (*domain.PriceResponse, error) {
	fullKey := c.prefix + key
	
	data, err := c.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var response domain.PriceResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	
	return &response, nil
}

// Set stores a price response in cache
func (c *RedisCache) Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error {
	fullKey := c.prefix + key
	
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, fullKey, data, ttl).Err()
}

// Invalidate removes cached entries matching the given patterns
func (c *RedisCache) Invalidate(ctx context.Context, patterns []string) error {
	for _, pattern := range patterns {
		fullPattern := c.prefix + pattern
		
		// Use SCAN to find matching keys
		iter := c.client.Scan(ctx, 0, fullPattern, 0).Iterator()
		keys := make([]string, 0)
		
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		
		if err := iter.Err(); err != nil {
			return err
		}
		
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// InvalidateForTenant invalidates all cache entries for a specific tenant
func (c *RedisCache) InvalidateForTenant(ctx context.Context, tenantID string) error {
	return c.Invalidate(ctx, []string{"price:" + tenantID + ":*"})
}

// InvalidateForCustomer invalidates all cache entries for a specific customer
func (c *RedisCache) InvalidateForCustomer(ctx context.Context, tenantID, customerID string) error {
	return c.Invalidate(ctx, []string{"price:" + tenantID + ":" + customerID + ":*"})
}

// InvalidateForProduct invalidates all cache entries containing a specific product
func (c *RedisCache) InvalidateForProduct(ctx context.Context, tenantID, productID string) error {
	return c.Invalidate(ctx, []string{"*:" + productID + ":*"})
}

// NullCache is a no-op cache implementation for testing or when caching is disabled
type NullCache struct{}

// NewNullCache creates a new null cache instance
func NewNullCache() *NullCache {
	return &NullCache{}
}

// Get always returns nil
func (c *NullCache) Get(ctx context.Context, key string) (*domain.PriceResponse, error) {
	return nil, nil
}

// Set always succeeds without storing anything
func (c *NullCache) Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error {
	return nil
}

// Invalidate always succeeds
func (c *NullCache) Invalidate(ctx context.Context, patterns []string) error {
	return nil
}

// InMemoryCache provides a simple in-memory cache for single-node deployments
type InMemoryCache struct {
	store     map[string]cacheEntry
	maxSize   int
	evictChan chan string
}

type cacheEntry struct {
	response  *domain.PriceResponse
	expiresAt time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache(maxSize int) *InMemoryCache {
	c := &InMemoryCache{
		store:     make(map[string]cacheEntry),
		maxSize:   maxSize,
		evictChan: make(chan string, 100),
	}
	
	// Start background cleanup goroutine
	go c.cleanup()
	
	return c
}

// Get retrieves a cached price response
func (c *InMemoryCache) Get(ctx context.Context, key string) (*domain.PriceResponse, error) {
	entry, ok := c.store[key]
	if !ok {
		return nil, nil
	}
	
	if time.Now().After(entry.expiresAt) {
		delete(c.store, key)
		return nil, nil
	}
	
	return entry.response, nil
}

// Set stores a price response in cache
func (c *InMemoryCache) Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error {
	// Simple eviction when at capacity
	if len(c.store) >= c.maxSize {
		// Evict one random entry
		for k := range c.store {
			delete(c.store, k)
			break
		}
	}
	
	c.store[key] = cacheEntry{
		response:  response,
		expiresAt: time.Now().Add(ttl),
	}
	
	return nil
}

// Invalidate removes cached entries matching patterns
func (c *InMemoryCache) Invalidate(ctx context.Context, patterns []string) error {
	// Simple implementation - clear all for now
	// In production, implement pattern matching
	for key := range c.store {
		delete(c.store, key)
	}
	return nil
}

// cleanup removes expired entries periodically
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		for key, entry := range c.store {
			if now.After(entry.expiresAt) {
				delete(c.store, key)
			}
		}
	}
}
