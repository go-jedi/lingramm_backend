package refreshtoken

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupCache() *RefreshToken {
	cfg := config.RedisConfig{
		Addr:            "127.0.0.1:63790",
		Password:        "auth",
		DB:              0,
		DialTimeout:     5,
		ReadTimeout:     3,
		WriteTimeout:    3,
		PoolSize:        10,
		MinIdleConns:    3,
		PoolTimeout:     4,
		MaxRetries:      3,
		MinRetryBackoff: 8,
		MaxRetryBackoff: 512,
		RefreshToken: config.RefreshTokenConfig{
			QueryTimeout: 2,
			Expiration:   7,
		},
	}

	c := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		PoolTimeout:     time.Duration(cfg.PoolTimeout) * time.Second,
		TLSConfig:       nil,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: time.Duration(cfg.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(cfg.MaxRetryBackoff) * time.Millisecond,
	})

	return New(cfg, c)
}

func TestSet(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
		val string
	}

	type want struct {
		refreshToken string
		err          error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx: ctx,
				key: telegramID,
				val: refreshToken,
			},
			want: want{
				refreshToken: refreshToken,
				err:          nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(ctx, test.in.key, test.in.val)
			assert.NoError(t, err)

			got, err := cache.Get(ctx, test.in.key)
			assert.Equal(t, test.want.refreshToken, got)
		})
	}
}

func TestAll(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
		val string
	}

	type want struct {
		key string
		val string
		err error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx: ctx,
				key: telegramID,
				val: refreshToken,
			},
			want: want{
				key: telegramID,
				val: refreshToken,
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(ctx, test.in.key, test.in.val)
			assert.NoError(t, err)

			result, err := cache.All(test.in.ctx)
			assert.NoError(t, err)

			res, ok := result[cache.getPrefixRefreshToken()+cache.getPrefixTelegramID()+test.in.key]
			assert.True(t, ok)
			assert.Equal(t, test.want.val, res)
		})
	}
}

func TestGet(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
		val string
	}

	type want struct {
		val string
		err error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx: ctx,
				key: telegramID,
				val: refreshToken,
			},
			want: want{
				val: refreshToken,
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(ctx, test.in.key, test.in.val)
			assert.NoError(t, err)

			result, err := cache.Get(ctx, test.in.key)
			assert.NoError(t, err)
			assert.Equal(t, test.want.val, result)
		})
	}
}

func TestExists(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
		val string
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx: ctx,
				key: telegramID,
				val: refreshToken,
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
		{
			name: "err",
			in: in{
				ctx: ctx,
				key: gofakeit.UUID(),
				val: refreshToken,
			},
			want: want{
				exists: false,
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			switch test.name {
			case "ok":
				err := cache.Set(ctx, test.in.key, test.in.val)
				assert.NoError(t, err)

				exists, err := cache.Exists(ctx, test.in.key)
				assert.NoError(t, err)
				assert.Equal(t, test.want.exists, exists)
			default:
				exists, err := cache.Exists(ctx, test.in.key)
				assert.NoError(t, err)
				assert.Equal(t, test.want.exists, exists)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
		val string
	}

	type want struct {
		err error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx: ctx,
				key: telegramID,
				val: refreshToken,
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(ctx, test.in.key, test.in.val)
			assert.NoError(t, err)

			err = cache.Delete(ctx, test.in.key)
			assert.NoError(t, err)
		})
	}
}

func TestDeleteKeys(t *testing.T) {
	type in struct {
		ctx  context.Context
		keys []string
	}

	type want struct {
		err error
	}

	var (
		ctx      = context.TODO()
		testKeys = []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
		}
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx:  ctx,
				keys: testKeys,
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(ctx, test.in.keys[0], gofakeit.UUID())
			assert.NoError(t, err)
			err = cache.Set(ctx, test.in.keys[1], gofakeit.UUID())
			assert.NoError(t, err)

			err = cache.DeleteKeys(ctx, test.in.keys)
			assert.NoError(t, err)
		})
	}
}
