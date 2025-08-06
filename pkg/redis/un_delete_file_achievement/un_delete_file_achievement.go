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
	prefixFileName                = "file_name:"
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
	prefixFileName                string
}

func New(cfg config.UnDeleteFileAchievementConfig, client *redis.Client) *UnDeleteFileAchievement {
	return &UnDeleteFileAchievement{
		client:                        client,
		prefixUnDeleteFileAchievement: prefixUnDeleteFileAchievement,
		prefixFileName:                prefixFileName,
		queryTimeout:                  cfg.QueryTimeout,
		expiration:                    cfg.Expiration,
	}
}

// Set stores un delete file achievement in Redis using MessagePack serialization.
func (c *UnDeleteFileAchievement) Set(ctx context.Context, key string, val string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctxTimeout,
		c.getRedisKey(key),
		b,
		c.getExpiration(),
	).Err()
}

// All retrieves all un delete files achievement entries from the cache.
func (c *UnDeleteFileAchievement) All(ctx context.Context) (map[string]string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	const count = 200
	var (
		cursor uint64
		result = make(map[string]string)
		match  = c.getPrefixUnDeleteFileAchievement() + c.getPrefixFileName() + "*"
	)

	for {
		keys, nextCursor, err := c.client.Scan(ctxTimeout, cursor, match, count).Result()
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

		values, err := c.client.MGet(ctxTimeout, keys...).Result()
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

// Delete removes un delete file achievement from the cache by key.
func (c *UnDeleteFileAchievement) Delete(ctx context.Context, key string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctxTimeout, key).Err()
}

// DeleteKeys removes un delete files achievement from the cache by keys.
func (c *UnDeleteFileAchievement) DeleteKeys(ctx context.Context, keys []string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctxTimeout, keys...).Err()
}

// getRedisKey get redis key.
func (c *UnDeleteFileAchievement) getRedisKey(key string) string {
	return c.getPrefixUnDeleteFileAchievement() + c.getPrefixFileName() + key
}

// getPrefixUnDeleteFileAchievement get prefix un delete file achievement.
func (c *UnDeleteFileAchievement) getPrefixUnDeleteFileAchievement() string {
	return c.prefixUnDeleteFileAchievement
}

// getPrefixFileName get prefix file name.
func (c *UnDeleteFileAchievement) getPrefixFileName() string {
	return c.prefixFileName
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
