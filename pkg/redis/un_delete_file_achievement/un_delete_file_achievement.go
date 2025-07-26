package undeletefileachievement

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixUnDeleteFileAchievement = "un_delete_file_achievement:"
	prefixFileID                  = "file_id:"
)

//go:generate mockery --name=IUnDeleteFileAchievement --output=mocks --case=underscore
type IUnDeleteFileAchievement interface {
	Set(ctx context.Context, key string, val string) error
	All(ctx context.Context) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
}

type UnDeleteFileAchievement struct {
	queryTimeout                  int64
	expiration                    int64
	client                        *redis.Client
	prefixUnDeleteFileAchievement string
	prefixFileID                  string
}

func New(cfg config.UnDeleteFileAchievementConfig, client *redis.Client) *UnDeleteFileAchievement {
	return &UnDeleteFileAchievement{
		client:                        client,
		prefixUnDeleteFileAchievement: prefixUnDeleteFileAchievement,
		prefixFileID:                  prefixFileID,
		queryTimeout:                  cfg.QueryTimeout,
		expiration:                    cfg.Expiration,
	}
}

// Set stores refresh token in Redis using MessagePack serialization.
func (c *UnDeleteFileAchievement) Set(ctx context.Context, key string, val string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctx,
		c.getPrefixUnDeleteFileAchievement()+c.getPrefixFileID()+key,
		b,
		c.getExpiration(),
	).Err()
}

// All retrieves all refresh token entries from the cache.
func (c *UnDeleteFileAchievement) All(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	const count = 200
	var (
		cursor uint64
		result = make(map[string]string)
		match  = c.getPrefixUnDeleteFileAchievement() + c.getPrefixFileID() + "*"
	)

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, match, count).Result()
		if err != nil {
			return nil, err
		}

		if len(keys) == 0 {
			if nextCursor == 0 {
				break
			}
			cursor = nextCursor
			continue
		}

		values, err := c.client.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for i := range values {
			rawData := c.convertToBytes(values[i])
			if rawData == nil {
				continue
			}

			var item string
			if err := msgpack.Unmarshal(rawData, &item); err != nil {
				continue
			}

			result[keys[i]] = item
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return result, nil
}

// Delete removes refresh token from the cache by key.
func (c *UnDeleteFileAchievement) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, key).Err()
}

// DeleteKeys removes refresh tokens from the cache by keys.
func (c *UnDeleteFileAchievement) DeleteKeys(ctx context.Context, keys []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, keys...).Err()
}

// getPrefixUnDeleteFileAchievement get prefix un delete file.
func (c *UnDeleteFileAchievement) getPrefixUnDeleteFileAchievement() string {
	return c.prefixUnDeleteFileAchievement
}

// getPrefixFileID get prefix file id.
func (c *UnDeleteFileAchievement) getPrefixFileID() string {
	return c.prefixFileID
}

// getExpiration get expiration date for row in cache.
func (c *UnDeleteFileAchievement) getExpiration() time.Duration {
	return time.Duration(c.expiration) * 24 * time.Hour
}

// convertToBytes safely converts interface{} to []byte.
func (c *UnDeleteFileAchievement) convertToBytes(val interface{}) []byte {
	switch v := val.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}
