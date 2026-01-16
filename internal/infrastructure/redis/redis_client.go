package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	SetIdempotencyKey(ctx context.Context, key string, value string, expiration time.Duration) error
	CheckIdempotencyKey(ctx context.Context, key string) (bool, error)
	AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
	SetCache(ctx context.Context, key string, value string, expiration time.Duration) error
	GetCache(ctx context.Context, key string) (string, error)
	Close() error
}

type Client struct {
	client *redis.Client
}

func NewClient(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{client: rdb}, nil
}

func (c *Client) SetIdempotencyKey(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.SetNX(ctx, fmt.Sprintf("idempotency:%s", key), value, expiration).Err()
}

func (c *Client) CheckIdempotencyKey(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, fmt.Sprintf("idempotency:%s", key)).Result()
	return exists > 0, err
}

func (c *Client) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, fmt.Sprintf("lock:%s", key), "1", expiration).Result()
}

func (c *Client) ReleaseLock(ctx context.Context, key string) error {
	return c.client.Del(ctx, fmt.Sprintf("lock:%s", key)).Err()
}

func (c *Client) SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, fmt.Sprintf("cache:%s", key), value, expiration).Err()
}

func (c *Client) GetCache(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, fmt.Sprintf("cache:%s", key)).Result()
}

func (c *Client) Close() error {
	return c.client.Close()
}
