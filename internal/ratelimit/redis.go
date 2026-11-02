// Package ratelimit provides shared, Redis-backed request limits.
package ratelimit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter checks and records a request against a named quota.
type Limiter interface {
	Allow(ctx context.Context, key string, limit int64, window time.Duration) (allowed bool, retryAfter time.Duration, err error)
}

type RedisLimiter struct {
	client *redis.Client
}

var fixedWindowScript = redis.NewScript(`
local count = redis.call('INCR', KEYS[1])
if count == 1 then
  redis.call('PEXPIRE', KEYS[1], ARGV[1])
end
return {count, redis.call('PTTL', KEYS[1])}
`)

// NewRedis creates a limiter from a redis:// or rediss:// connection URL.
// The connection is lazy: Redis outages only reject rate-limited operations,
// leaving the rest of the application available.
func NewRedis(rawURL string) (*RedisLimiter, error) {
	options, err := redis.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse Redis URL: %w", err)
	}
	return &RedisLimiter{client: redis.NewClient(options)}, nil
}

func (l *RedisLimiter) Close() error {
	return l.client.Close()
}

func (l *RedisLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, time.Duration, error) {
	if limit < 1 || window <= 0 {
		return false, 0, fmt.Errorf("rate limit must have a positive limit and window")
	}
	digest := sha256.Sum256([]byte(key))
	redisKey := "rate-limit:" + hex.EncodeToString(digest[:])

	result, err := fixedWindowScript.Run(ctx, l.client, []string{redisKey}, window.Milliseconds()).Slice()
	if err != nil {
		return false, 0, fmt.Errorf("check Redis rate limit: %w", err)
	}
	if len(result) != 2 {
		return false, 0, fmt.Errorf("unexpected Redis rate limit response")
	}
	count, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected Redis rate limit counter type %T", result[0])
	}
	ttl, ok := result[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected Redis rate limit TTL type %T", result[1])
	}
	return count <= limit, time.Duration(ttl) * time.Millisecond, nil
}
