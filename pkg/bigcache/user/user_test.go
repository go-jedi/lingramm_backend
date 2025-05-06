package user

import (
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func setupCache(t *testing.T) *User {
	config := bigcache.DefaultConfig(10 * time.Minute)
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		t.Fatalf("failed to create bigcache: %v", err)
	}
	return New(cache)
}

func TestSet(t *testing.T) {
	type in struct {
		key string
		val user.User
	}

	type want struct {
		user user.User
		err  error
	}

	var (
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
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
				key: telegramID,
				val: testUser,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			err := cache.Set(test.in.key, test.in.val)
			assert.NoError(t, err)

			got, err := cache.Get(test.in.key)
			assert.Equal(t, test.want.err, err)

			assert.Equal(t, test.want.user.ID, got.ID)
			assert.Equal(t, test.want.user.UUID, got.UUID)
			assert.Equal(t, test.want.user.TelegramID, got.TelegramID)
			assert.Equal(t, test.want.user.Username, got.Username)
			assert.Equal(t, test.want.user.FirstName, got.FirstName)
			assert.Equal(t, test.want.user.LastName, got.LastName)

			assert.WithinDuration(t, test.want.user.CreatedAt, got.CreatedAt, time.Millisecond)
			assert.WithinDuration(t, test.want.user.UpdatedAt, got.UpdatedAt, time.Millisecond)
		})
	}
}

func TestAll(t *testing.T) {
	type in struct {
		key string
		val user.User
	}

	type want struct {
		user []user.User
		err  error
	}

	var (
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
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
				key: telegramID,
				val: testUser,
			},
			want: want{
				user: []user.User{testUser},
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			err := cache.Set(test.in.key, test.in.val)
			assert.NoError(t, err)

			got, err := cache.All()
			assert.Equal(t, test.want.err, err)

			assert.Len(t, got, len(test.want.user))

			assert.Equal(t, test.want.user[0].ID, got[0].ID)
			assert.Equal(t, test.want.user[0].UUID, got[0].UUID)
			assert.Equal(t, test.want.user[0].TelegramID, got[0].TelegramID)
			assert.Equal(t, test.want.user[0].Username, got[0].Username)
			assert.Equal(t, test.want.user[0].FirstName, got[0].FirstName)
			assert.Equal(t, test.want.user[0].LastName, got[0].LastName)

			assert.WithinDuration(t, test.want.user[0].CreatedAt, got[0].CreatedAt, time.Millisecond)
			assert.WithinDuration(t, test.want.user[0].UpdatedAt, got[0].UpdatedAt, time.Millisecond)
		})
	}
}

func TestGet(t *testing.T) {
	type in struct {
		key string
		val user.User
	}

	type want struct {
		user user.User
		err  error
	}

	var (
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
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
				key: telegramID,
				val: testUser,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
		{
			name: "not found",
			in: in{
				key: telegramID,
				val: user.User{},
			},
			want: want{
				user: user.User{},
				err:  bigcache.ErrEntryNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			switch test.name {
			case "ok":
				err := cache.Set(test.in.key, test.in.val)
				assert.NoError(t, err)

				got, err := cache.Get(test.in.key)
				assert.Equal(t, test.want.err, err)

				assert.Equal(t, test.want.user.ID, got.ID)
				assert.Equal(t, test.want.user.UUID, got.UUID)
				assert.Equal(t, test.want.user.TelegramID, got.TelegramID)
				assert.Equal(t, test.want.user.Username, got.Username)
				assert.Equal(t, test.want.user.FirstName, got.FirstName)
				assert.Equal(t, test.want.user.LastName, got.LastName)

				assert.WithinDuration(t, test.want.user.CreatedAt, got.CreatedAt, time.Millisecond)
				assert.WithinDuration(t, test.want.user.UpdatedAt, got.UpdatedAt, time.Millisecond)
			default:
				_, err := cache.Get(test.in.key)
				assert.Equal(t, test.want.err, err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type in struct {
		key string
		val user.User
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
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
				key: telegramID,
				val: testUser,
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
		{
			name: "not found",
			in: in{
				key: telegramID,
				val: testUser,
			},
			want: want{
				exists: false,
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			switch test.name {
			case "ok":
				err := cache.Set(test.in.key, test.in.val)
				assert.NoError(t, err)

				got, err := cache.Exists(test.in.key)
				assert.Equal(t, test.want.err, err)
				assert.Equal(t, test.want.exists, got)
			default:
				got, err := cache.Exists(test.in.key)
				assert.Equal(t, test.want.err, err)
				assert.Equal(t, test.want.exists, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type in struct {
		key string
		val user.User
	}

	type want struct {
		user user.User
		err  error
	}

	var (
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
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
				key: telegramID,
				val: testUser,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := setupCache(t)

			err := cache.Set(test.in.key, test.in.val)
			assert.NoError(t, err)

			err = cache.Delete(test.in.key)
			assert.NoError(t, err)
		})
	}
}
