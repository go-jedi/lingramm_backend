package localizedtext

import (
	"testing"
	"time"

	"github.com/allegro/bigcache"
	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	"github.com/stretchr/testify/assert"
)

func setupCache(t *testing.T) *LocalizedText {
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
		val map[string][]localizedtext.LocalizedTexts
	}

	type want struct {
		localizedTexts map[string][]localizedtext.LocalizedTexts
		err            error
	}

	var (
		key                = "en"
		description        = "Добро пожаловать!"
		testLocalizedTexts = map[string][]localizedtext.LocalizedTexts{
			key: {
				{
					Code:        "welcome_title",
					Value:       "Добро пожаловать!",
					Description: &description,
				},
			},
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
				key: key,
				val: testLocalizedTexts,
			},
			want: want{
				localizedTexts: testLocalizedTexts,
				err:            nil,
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

			assert.Equal(t, test.want.localizedTexts[test.in.key][0].Code, got[test.in.key][0].Code)
			assert.Equal(t, test.want.localizedTexts[test.in.key][0].Value, got[test.in.key][0].Value)
			assert.Equal(t, test.want.localizedTexts[test.in.key][0].Description, got[test.in.key][0].Description)
		})
	}
}

func TestGet(t *testing.T) {
	type in struct {
		key string
		val map[string][]localizedtext.LocalizedTexts
	}

	type want struct {
		localizedTexts map[string][]localizedtext.LocalizedTexts
		err            error
	}

	var (
		key                = "en"
		description        = "Добро пожаловать!"
		testLocalizedTexts = map[string][]localizedtext.LocalizedTexts{
			key: {
				{
					Code:        "welcome_title",
					Value:       "Добро пожаловать!",
					Description: &description,
				},
			},
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
				key: key,
				val: testLocalizedTexts,
			},
			want: want{
				localizedTexts: testLocalizedTexts,
				err:            nil,
			},
		},
		{
			name: "not found",
			in: in{
				key: key,
				val: nil,
			},
			want: want{
				localizedTexts: nil,
				err:            bigcache.ErrEntryNotFound,
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

				assert.Equal(t, test.want.localizedTexts[test.in.key][0].Code, got[test.in.key][0].Code)
				assert.Equal(t, test.want.localizedTexts[test.in.key][0].Value, got[test.in.key][0].Value)
				assert.Equal(t, test.want.localizedTexts[test.in.key][0].Description, got[test.in.key][0].Description)
			default:
				_, err := cache.Get(test.in.key)
				assert.Equal(t, test.want.err, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type in struct {
		key string
		val map[string][]localizedtext.LocalizedTexts
	}

	type want struct {
		localizedTexts map[string][]localizedtext.LocalizedTexts
		err            error
	}

	var (
		key                = "en"
		description        = "Добро пожаловать!"
		testLocalizedTexts = map[string][]localizedtext.LocalizedTexts{
			key: {
				{
					Code:        "welcome_title",
					Value:       "Добро пожаловать!",
					Description: &description,
				},
			},
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
				key: key,
				val: testLocalizedTexts,
			},
			want: want{
				localizedTexts: testLocalizedTexts,
				err:            nil,
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
