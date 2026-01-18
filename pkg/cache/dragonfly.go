// Package cache provides DragonflyDB (Redis-compatible) client configuration.
// DragonflyDB offers 25x better memory efficiency than Redis with full API compatibility.
package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DragonflyConfig holds DragonflyDB connection configuration
type DragonflyConfig struct {
	// Addresses for cluster mode
	Addresses []string
	// Password for authentication
	Password string
	// DB number (0-15)
	DB int
	// PoolSize is the maximum number of connections
	PoolSize int
	// MinIdleConns is the minimum number of idle connections
	MinIdleConns int
	// MaxRetries before failing
	MaxRetries int
	// TLSEnabled enables TLS connections
	TLSEnabled bool
	// DialTimeout for connections
	DialTimeout time.Duration
	// ReadTimeout for reads
	ReadTimeout time.Duration
	// WriteTimeout for writes
	WriteTimeout time.Duration
}

// DefaultDragonflyConfig returns default configuration
func DefaultDragonflyConfig() DragonflyConfig {
	return DragonflyConfig{
		Addresses:    []string{"localhost:6379"},
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// NewDragonflyClient creates a new DragonflyDB client (single node)
func NewDragonflyClient(cfg DragonflyConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.Addresses[0],
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	if cfg.TLSEnabled {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opts)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return client, nil
}

// NewDragonflyClusterClient creates a new DragonflyDB cluster client
func NewDragonflyClusterClient(cfg DragonflyConfig) (*redis.ClusterClient, error) {
	opts := &redis.ClusterOptions{
		Addrs:         cfg.Addresses,
		Password:      cfg.Password,
		PoolSize:      cfg.PoolSize,
		MinIdleConns:  cfg.MinIdleConns,
		MaxRetries:    cfg.MaxRetries,
		DialTimeout:   cfg.DialTimeout,
		ReadTimeout:   cfg.ReadTimeout,
		WriteTimeout:  cfg.WriteTimeout,
		RouteRandomly: true,
	}

	if cfg.TLSEnabled {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClusterClient(opts)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return client, nil
}

// CacheService provides high-level caching operations
type CacheService struct {
	client redis.UniversalClient
}

// NewCacheService creates a new cache service
func NewCacheService(client redis.UniversalClient) *CacheService {
	return &CacheService{client: client}
}

// Set stores a value with TTL
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key
func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// SetNX sets a value only if it doesn't exist (for distributed locking)
func (c *CacheService) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("marshal: %w", err)
	}
	return c.client.SetNX(ctx, key, data, ttl).Result()
}

// AcquireLock attempts to acquire a distributed lock
func (c *CacheService) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

// ReleaseLock releases a distributed lock
func (c *CacheService) ReleaseLock(ctx context.Context, key string) error {
	return c.client.Del(ctx, "lock:"+key).Err()
}

// Incr increments a counter
func (c *CacheService) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// CheckRateLimit checks if a rate limit is exceeded
func (c *CacheService) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	pipe := c.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return incr.Val() <= int64(limit), nil
}

// HSet sets hash fields
func (c *CacheService) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return c.client.HSet(ctx, key, values).Err()
}

// HGetAll gets all hash fields
func (c *CacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// Publish publishes a message to a channel
func (c *CacheService) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return c.client.Publish(ctx, channel, data).Err()
}

// Subscribe subscribes to channels
func (c *CacheService) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}
