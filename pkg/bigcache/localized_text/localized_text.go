package localizedtext

import (
	"errors"

	"github.com/allegro/bigcache"
	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixLocalizedText = "localized_text:"
	prefixLanguage      = "language:"
)

//go:generate mockery --name=ILocalizedText --output=mocks --case=underscore
type ILocalizedText interface {
	Set(key string, val map[string][]localizedtext.LocalizedTexts) error
	Get(key string) (map[string][]localizedtext.LocalizedTexts, error)
	Delete(key string) error
}

type LocalizedText struct {
	prefixLocalizedText string
	prefixLanguage      string
	bigCache            *bigcache.BigCache
}

func New(bigCache *bigcache.BigCache) *LocalizedText {
	return &LocalizedText{
		prefixLocalizedText: prefixLocalizedText,
		prefixLanguage:      prefixLanguage,
		bigCache:            bigCache,
	}
}

// Set stores localized text in BigCache using MessagePack serialization.
func (c *LocalizedText) Set(key string, val map[string][]localizedtext.LocalizedTexts) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return c.bigCache.Set(c.getPrefixLocalizedText()+c.getPrefixLanguage()+key, b)
}

// Get retrieves localized text from BigCache and deserializes it using MessagePack.
func (c *LocalizedText) Get(key string) (map[string][]localizedtext.LocalizedTexts, error) {
	var result map[string][]localizedtext.LocalizedTexts

	data, err := c.bigCache.Get(c.getPrefixLocalizedText() + c.getPrefixLanguage() + key)
	if err != nil {
		return nil, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes localized text from the cache by key.
func (c *LocalizedText) Delete(key string) error {
	err := c.bigCache.Delete(c.getPrefixLocalizedText() + c.getPrefixLanguage() + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil
		}
		return err
	}

	return nil
}

// getPrefixLocalizedText get prefix localized text.
func (c *LocalizedText) getPrefixLocalizedText() string {
	return c.prefixLocalizedText
}

// getPrefixLanguage get prefix language.
func (c *LocalizedText) getPrefixLanguage() string {
	return c.prefixLanguage
}
