package refreshtoken

import (
	"context"
	"errors"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixRefreshToken = "refresh_token:"
	prefixTelegramID   = "telegram_id:"
)

//go:generate mockery --name=IRefreshToken --output=mocks --case=underscore
type IRefreshToken interface {
	Set(ctx context.Context, key string, val string) error
	SetWithExpiration(ctx context.Context, key string, val string, expiration time.Duration) error
	All(ctx context.Context) (map[string]string, error)
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
}

type RefreshToken struct {
	queryTimeout       int64
	expiration         int64
	client             *redis.Client
	prefixRefreshToken string
	prefixTelegramID   string
}

func New(cfg config.RedisConfig, client *redis.Client) *RefreshToken {
	return &RefreshToken{
		client:             client,
		prefixRefreshToken: prefixRefreshToken,
		prefixTelegramID:   prefixTelegramID,
		queryTimeout:       cfg.RefreshToken.QueryTimeout,
		expiration:         cfg.RefreshToken.Expiration,
	}
}

// Set stores refresh token in Redis using MessagePack serialization.
func (rf *RefreshToken) Set(ctx context.Context, key string, val string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	return rf.client.Set(
		ctx,
		rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+key,
		b,
		rf.getExpiration(),
	).Err()
}

// SetWithExpiration set stores refresh token with expiration in Redis using MessagePack serialization.
func (rf *RefreshToken) SetWithExpiration(ctx context.Context, key string, val string, expiration time.Duration) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	return rf.client.Set(
		ctx,
		rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+key,
		b,
		expiration,
	).Err()
}

// All retrieves all refresh token entries from the cache.
func (rf *RefreshToken) All(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	const count = 200
	var (
		cursor uint64
		result = make(map[string]string)
		match  = rf.getPrefixRefreshToken() + rf.getPrefixTelegramID() + "*"
	)

	for {
		keys, nextCursor, err := rf.client.Scan(ctx, cursor, match, count).Result()
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

		values, err := rf.client.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for i := range values {
			rawData := rf.convertToBytes(values[i])
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

// Get retrieves refresh token from Redis and deserializes it using MessagePack.
func (rf *RefreshToken) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	b, err := rf.client.Get(
		ctx,
		rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+key,
	).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}

	var result string
	if err := msgpack.Unmarshal(b, &result); err != nil {
		return "", err
	}

	return result, nil
}

// Exists checks whether a refresh token exists in Redis by key.
func (rf *RefreshToken) Exists(ctx context.Context, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	count, err := rf.client.Exists(ctx, rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete removes refresh token from the cache by key.
func (rf *RefreshToken) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	return rf.client.Del(ctx, rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+key).Err()
}

// DeleteKeys removes refresh tokens from the cache by keys.
func (rf *RefreshToken) DeleteKeys(ctx context.Context, keys []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rf.queryTimeout)*time.Second)
	defer cancel()

	fullKeys := make([]string, 0, len(keys))
	for i := range keys {
		fullKeys = append(fullKeys, rf.getPrefixRefreshToken()+rf.getPrefixTelegramID()+keys[i])
	}

	return rf.client.Del(ctx, fullKeys...).Err()
}

// getPrefixRefreshToken get prefix refresh token.
func (rf *RefreshToken) getPrefixRefreshToken() string {
	return rf.prefixRefreshToken
}

// getPrefixTelegramID get prefix telegram id.
func (rf *RefreshToken) getPrefixTelegramID() string {
	return rf.prefixTelegramID
}

// getExpiration get expiration date for row in cache.
func (rf *RefreshToken) getExpiration() time.Duration {
	return time.Duration(rf.expiration) * 24 * time.Hour
}

// convertToBytes safely converts interface{} to []byte.
func (rf *RefreshToken) convertToBytes(val interface{}) []byte {
	switch v := val.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}
