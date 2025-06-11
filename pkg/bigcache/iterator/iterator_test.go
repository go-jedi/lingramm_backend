package iterator

import (
	"context"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func setupCache(t *testing.T) *bigcache.BigCache {
	config := bigcache.DefaultConfig(10 * time.Minute)

	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		t.Fatalf("failed to create bigcache: %v", err)
	}

	return cache
}

func setDataInCache(key string, value any, prefix string, bigCache *bigcache.BigCache) error {
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	if err := bigCache.Set(prefix+key, bytes); err != nil {
		return err
	}

	return nil
}

func TestIterator(t *testing.T) {
	type in struct {
		ctx    context.Context
		key    string
		val    interface{}
		prefix string
	}

	type want struct {
		val user.User
		err error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			Username:   gofakeit.Username(),
			FirstName:  gofakeit.FirstName(),
			LastName:   gofakeit.LastName(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
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
				ctx:    ctx,
				key:    telegramID,
				val:    testUser,
				prefix: "user:telegram_id:",
			},
			want: want{
				val: testUser,
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			i := New(cache)

			err := setDataInCache(test.in.key, test.in.val, test.in.prefix, cache)
			assert.NoError(t, err)

			result, err := i.Iterator(test.in.ctx)
			assert.NoError(t, err)
			assert.Equal(t, test.want.err, err)

			out, ok := result[test.in.prefix+test.in.key].(*user.User)
			assert.True(t, ok)

			assert.Equal(t, test.want.val.ID, out.ID)
			assert.Equal(t, test.want.val.TelegramID, out.TelegramID)
			assert.Equal(t, test.want.val.Username, out.Username)
			assert.Equal(t, test.want.val.FirstName, out.FirstName)
			assert.Equal(t, test.want.val.LastName, out.LastName)
		})
	}
}

func TestRegisterType(t *testing.T) {
	type in struct {
		prefix  string
		factory TypeFactory
	}

	type want struct {
		key string
		ok  bool
	}

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "ok",
			in: in{
				prefix:  "user:telegram_id:",
				factory: func() any { return new(user.User) },
			},
			want: want{
				key: "user:telegram_id:",
				ok:  true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			i := New(cache)

			i.RegisterType(test.in.prefix, test.in.factory)

			_, ok := i.typeLookup[test.want.key]
			assert.Equal(t, test.want.ok, ok)
		})
	}
}
