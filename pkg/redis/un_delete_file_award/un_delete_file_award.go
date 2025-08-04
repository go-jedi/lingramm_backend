package undeletefileaward

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixUnDeleteFileAward = "un_delete_file_award:"
	prefixFileName          = "file_name:"
)

//go:generate mockery --name=IUnDeleteFileAward --output=mocks --case=underscore
type IUnDeleteFileAward interface {
	Set(ctx context.Context, key string, val string) error
	All(ctx context.Context) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
}

type UnDeleteFileAward struct {
	queryTimeout            int64
	expiration              int64
	client                  *redis.Client
	prefixUnDeleteFileAward string
	prefixFileName          string
}

func New(cfg config.UnDeleteFileAwardConfig, client *redis.Client) *UnDeleteFileAward {
	return &UnDeleteFileAward{
		client:                  client,
		prefixUnDeleteFileAward: prefixUnDeleteFileAward,
		prefixFileName:          prefixFileName,
		queryTimeout:            cfg.QueryTimeout,
		expiration:              cfg.Expiration,
	}
}

// Set stores un delete file award in Redis using MessagePack serialization.
func (c *UnDeleteFileAward) Set(ctx context.Context, key string, val string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctx,
		c.getPrefixUnDeleteFileAward()+c.getPrefixFileName()+key,
		b,
		c.getExpiration(),
	).Err()
}

// All retrieves all un delete files award entries from the cache.
func (c *UnDeleteFileAward) All(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	const count = 200
	var (
		cursor uint64
		result = make(map[string]string)
		match  = c.getPrefixUnDeleteFileAward() + c.getPrefixFileName() + "*"
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

// Delete removes un delete file award from the cache by key.
func (c *UnDeleteFileAward) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, key).Err()
}

// DeleteKeys removes un delete files award from the cache by keys.
func (c *UnDeleteFileAward) DeleteKeys(ctx context.Context, keys []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, keys...).Err()
}

// getPrefixUnDeleteFileAward get prefix un delete file award.
func (c *UnDeleteFileAward) getPrefixUnDeleteFileAward() string {
	return c.prefixUnDeleteFileAward
}

// getPrefixFileName get prefix file name.
func (c *UnDeleteFileAward) getPrefixFileName() string {
	return c.prefixFileName
}

// getExpiration get expiration date for row in cache.
func (c *UnDeleteFileAward) getExpiration() time.Duration {
	return time.Duration(c.expiration) * 24 * time.Hour
}

// convertToBytes safely converts interface{} to []byte.
func (c *UnDeleteFileAward) convertToBytes(val interface{}) []byte {
	switch v := val.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}
