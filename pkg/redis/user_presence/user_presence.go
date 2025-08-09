package userpresence

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
)

const (
	prefixUserPresence = "user_presence:"
	prefixTelegramID   = "telegram_id:"
)

//go:generate mockery --name=IUserPresence --output=mocks --case=underscore
type IUserPresence interface {
	Set(ctx context.Context, key string) error
	SetWithExpiration(ctx context.Context, key string, expiration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
}

type UserPresence struct {
	queryTimeout       int64
	expiration         int64
	client             *redis.Client
	prefixUserPresence string
	prefixTelegramID   string
}

func New(cfg config.UserPresenceConfig, client *redis.Client) *UserPresence {
	return &UserPresence{
		client:             client,
		prefixUserPresence: prefixUserPresence,
		prefixTelegramID:   prefixTelegramID,
		queryTimeout:       cfg.QueryTimeout,
		expiration:         cfg.Expiration,
	}
}

// Set stores user presence in Redis.
func (c *UserPresence) Set(ctx context.Context, key string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctxTimeout,
		c.getRedisKey(key),
		nil,
		c.getExpiration(),
	).Err()
}

// SetWithExpiration set stores user presence with expiration in Redis.
func (c *UserPresence) SetWithExpiration(ctx context.Context, key string, expiration time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctxTimeout,
		c.getRedisKey(key),
		nil,
		expiration,
	).Err()
}

// Exists checks whether a user presence exists in Redis by key.
func (c *UserPresence) Exists(ctx context.Context, key string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	count, err := c.client.Exists(ctxTimeout, c.getRedisKey(key)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete removes user presence from the cache by key.
func (c *UserPresence) Delete(ctx context.Context, key string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctxTimeout, key).Err()
}

// DeleteKeys removes users presences from the cache by keys.
func (c *UserPresence) DeleteKeys(ctx context.Context, keys []string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctxTimeout, keys...).Err()
}

// getRedisKey get redis key.
func (c *UserPresence) getRedisKey(key string) string {
	return c.getPrefixUserPresence() + c.getPrefixTelegramID() + key
}

// getPrefixUserPresence get prefix user presence.
func (c *UserPresence) getPrefixUserPresence() string {
	return c.prefixUserPresence
}

// getPrefixTelegramID get prefix telegram id.
func (c *UserPresence) getPrefixTelegramID() string {
	return c.prefixTelegramID
}

// getExpiration get expiration date for row in cache.
func (c *UserPresence) getExpiration() time.Duration {
	return time.Duration(c.expiration) * time.Second
}
