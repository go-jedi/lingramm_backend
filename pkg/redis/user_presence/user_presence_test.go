package userpresence

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupCache() *UserPresence {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:63790"
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		password = "admin"
	}

	cfg := config.RedisConfig{
		Addr:            addr,
		Password:        password,
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
		UserPresence: config.UserPresenceConfig{
			QueryTimeout: 2,
			Expiration:   60,
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

	return New(cfg.UserPresence, c)
}

func TestSet(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
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
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(test.in.ctx, test.in.key)
			assert.NoError(t, err)

			got, err := cache.Exists(test.in.ctx, telegramID)
			assert.NoError(t, err)
			assert.Equal(t, test.want.exists, got)
		})
	}
}

func TestSetWithExpiration(t *testing.T) {
	type in struct {
		ctx        context.Context
		key        string
		expiration time.Duration
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				ctx:        ctx,
				key:        telegramID,
				expiration: time.Duration(60) * time.Second,
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.SetWithExpiration(test.in.ctx, test.in.key, test.in.expiration)
			assert.NoError(t, err)

			got, err := cache.Exists(test.in.ctx, test.in.key)
			assert.NoError(t, err)
			assert.Equal(t, test.want.exists, got)
		})
	}
}

func TestExists(t *testing.T) {
	type in struct {
		ctx context.Context
		key string
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
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
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
		{
			name: "error",
			in: in{
				ctx: ctx,
				key: gofakeit.UUID(),
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
				err := cache.Set(test.in.ctx, test.in.key)
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
	}

	type want struct {
		err error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
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
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache()

			err := cache.Set(test.in.ctx, test.in.key)
			assert.NoError(t, err)

			err = cache.Delete(test.in.ctx, test.in.key)
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

			err := cache.Set(ctx, test.in.keys[0])
			assert.NoError(t, err)
			err = cache.Set(ctx, test.in.keys[1])
			assert.NoError(t, err)

			err = cache.DeleteKeys(ctx, test.in.keys)
			assert.NoError(t, err)
		})
	}
}
