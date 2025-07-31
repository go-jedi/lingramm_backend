package admin

import (
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	"github.com/stretchr/testify/assert"
)

func setupCache(t *testing.T) *Admin {
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
		val admin.Admin
	}

	type want struct {
		admin admin.Admin
		err   error
	}

	var (
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
				val: testAdmin,
			},
			want: want{
				admin: testAdmin,
				err:   nil,
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

			assert.Equal(t, test.want.admin.ID, got.ID)
			assert.Equal(t, test.want.admin.TelegramID, got.TelegramID)

			assert.WithinDuration(t, test.want.admin.CreatedAt, got.CreatedAt, time.Millisecond)
			assert.WithinDuration(t, test.want.admin.UpdatedAt, got.UpdatedAt, time.Millisecond)
		})
	}
}

func TestAll(t *testing.T) {
	type in struct {
		key string
		val admin.Admin
	}

	type want struct {
		admin []admin.Admin
		err   error
	}

	var (
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
				val: testAdmin,
			},
			want: want{
				admin: []admin.Admin{testAdmin},
				err:   nil,
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

			assert.Len(t, got, len(test.want.admin))

			assert.Equal(t, test.want.admin[0].ID, got[0].ID)
			assert.Equal(t, test.want.admin[0].TelegramID, got[0].TelegramID)

			assert.WithinDuration(t, test.want.admin[0].CreatedAt, got[0].CreatedAt, time.Millisecond)
			assert.WithinDuration(t, test.want.admin[0].UpdatedAt, got[0].UpdatedAt, time.Millisecond)
		})
	}
}

func TestGet(t *testing.T) {
	type in struct {
		key string
		val admin.Admin
	}

	type want struct {
		admin admin.Admin
		err   error
	}

	var (
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
				val: testAdmin,
			},
			want: want{
				admin: testAdmin,
				err:   nil,
			},
		},
		{
			name: "not found",
			in: in{
				key: telegramID,
				val: admin.Admin{},
			},
			want: want{
				admin: admin.Admin{},
				err:   bigcache.ErrEntryNotFound,
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

				assert.Equal(t, test.want.admin.ID, got.ID)
				assert.Equal(t, test.want.admin.TelegramID, got.TelegramID)

				assert.WithinDuration(t, test.want.admin.CreatedAt, got.CreatedAt, time.Millisecond)
				assert.WithinDuration(t, test.want.admin.UpdatedAt, got.UpdatedAt, time.Millisecond)
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
		val admin.Admin
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
				val: testAdmin,
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
				val: testAdmin,
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
		val admin.Admin
	}

	type want struct {
		admin admin.Admin
		err   error
	}

	var (
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
				val: testAdmin,
			},
			want: want{
				admin: testAdmin,
				err:   nil,
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
