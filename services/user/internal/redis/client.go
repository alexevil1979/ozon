package redis

import (
	"context"
	"fmt"

	redispkg "github.com/redis/go-redis/v9"
)

// Client is go-redis Client. Exported so main/handler can use without importing go-redis.
type Client = redispkg.Client

// NewClient creates Redis client from URL (e.g. redis://localhost:6379).
func NewClient(redisURL string) (*Client, error) {
	opt, err := redispkg.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	client := redispkg.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
