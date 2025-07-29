package undeletefileclient

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixUnDeleteFileClient = "un_delete_file_client:"
	prefixFileID             = "file_id:"
)

//go:generate mockery --name=IUnDeleteFileClient --output=mocks --case=underscore
type IUnDeleteFileClient interface {
	Set(ctx context.Context, key string, val string) error
	All(ctx context.Context) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
}

type UnDeleteFileClient struct {
	queryTimeout             int64
	expiration               int64
	client                   *redis.Client
	prefixUnDeleteFileClient string
	prefixFileID             string
}

func New(cfg config.UnDeleteFileClientConfig, client *redis.Client) *UnDeleteFileClient {
	return &UnDeleteFileClient{
		client:                   client,
		prefixUnDeleteFileClient: prefixUnDeleteFileClient,
		prefixFileID:             prefixFileID,
		queryTimeout:             cfg.QueryTimeout,
		expiration:               cfg.Expiration,
	}
}

// Set stores un delete file client in Redis using MessagePack serialization.
func (c *UnDeleteFileClient) Set(ctx context.Context, key string, val string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Set(
		ctx,
		c.getPrefixUnDeleteFileClient()+c.getPrefixFileID()+key,
		b,
		c.getExpiration(),
	).Err()
}

// All retrieves all un delete files client entries from the cache.
func (c *UnDeleteFileClient) All(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	const count = 200
	var (
		cursor uint64
		result = make(map[string]string)
		match  = c.getPrefixUnDeleteFileClient() + c.getPrefixFileID() + "*"
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

// Delete removes un delete file client from the cache by key.
func (c *UnDeleteFileClient) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, key).Err()
}

// DeleteKeys removes un delete files client from the cache by keys.
func (c *UnDeleteFileClient) DeleteKeys(ctx context.Context, keys []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	return c.client.Del(ctx, keys...).Err()
}

// getPrefixUnDeleteFileClient get prefix un delete file.
func (c *UnDeleteFileClient) getPrefixUnDeleteFileClient() string {
	return c.prefixUnDeleteFileClient
}

// getPrefixFileID get prefix file id.
func (c *UnDeleteFileClient) getPrefixFileID() string {
	return c.prefixFileID
}

// getExpiration get expiration date for row in cache.
func (c *UnDeleteFileClient) getExpiration() time.Duration {
	return time.Duration(c.expiration) * 24 * time.Hour
}

// convertToBytes safely converts interface{} to []byte.
func (c *UnDeleteFileClient) convertToBytes(val interface{}) []byte {
	switch v := val.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}
