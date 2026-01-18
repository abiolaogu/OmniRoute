// Package cache provides caching implementations for the pricing engine
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/omniroute/pricing-engine/internal/domain"
)

// RedisCache implements PriceCache using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{client: client, prefix: prefix}
}

// Get retrieves a cached price response
func (c *RedisCache) Get(ctx context.Context, key string) (*domain.PriceResponse, error) {
	data, err := c.client.Get(ctx, c.prefix+key).Bytes()
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

// Set stores a price response in the cache
func (c *RedisCache) Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.prefix+key, data, ttl).Err()
}

// Invalidate removes cached entries matching the given patterns
func (c *RedisCache) Invalidate(ctx context.Context, patterns []string) error {
	for _, pattern := range patterns {
		keys, err := c.client.Keys(ctx, c.prefix+pattern).Result()
		if err != nil {
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

// NullCache is a no-op cache implementation for when caching is disabled
type NullCache struct{}

// NewNullCache creates a new null cache
func NewNullCache() *NullCache {
	return &NullCache{}
}

// Get always returns nil (cache miss)
func (c *NullCache) Get(ctx context.Context, key string) (*domain.PriceResponse, error) {
	return nil, nil
}

// Set is a no-op
func (c *NullCache) Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error {
	return nil
}

// Invalidate is a no-op
func (c *NullCache) Invalidate(ctx context.Context, patterns []string) error {
	return nil
}
