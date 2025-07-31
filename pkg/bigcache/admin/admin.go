package admin

import (
	"errors"
	"log"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixAdmin      = "admin:"
	prefixTelegramID = "telegram_id:"
)

//go:generate mockery --name=IAdmin --output=mocks --case=underscore
type IAdmin interface {
	Set(key string, val admin.Admin) error
	All() ([]admin.Admin, error)
	Get(key string) (admin.Admin, error)
	Exists(key string) (bool, error)
	Delete(key string) error
}

type Admin struct {
	prefixAdmin      string
	prefixTelegramID string
	bigCache         *bigcache.BigCache
}

func New(bigCache *bigcache.BigCache) *Admin {
	return &Admin{
		prefixAdmin:      prefixAdmin,
		prefixTelegramID: prefixTelegramID,
		bigCache:         bigCache,
	}
}

// Set stores admin in BigCache using MessagePack serialization.
func (c *Admin) Set(key string, val admin.Admin) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return c.bigCache.Set(c.getPrefixAdmin()+c.getPrefixTelegramID()+key, b)
}

// All retrieves all admins entries from the cache that match the prefix.
func (c *Admin) All() ([]admin.Admin, error) {
	var result []admin.Admin

	iter := c.bigCache.Iterator()
	for iter.SetNext() {
		entry, err := iter.Value()
		if err != nil {
			log.Printf("failed to read cache entry: %v", err)
			continue
		}

		key := entry.Key()

		if !strings.HasPrefix(key, c.getPrefixAdmin()+c.getPrefixTelegramID()) {
			continue
		}

		var adm admin.Admin
		if err := msgpack.Unmarshal(entry.Value(), &adm); err != nil {
			log.Printf("failed to unmarshal cache entry for key '%s': %v", key, err)
			continue
		}

		result = append(result, adm)
	}

	return result, nil
}

// Get retrieves admin from BigCache and deserializes it using MessagePack.
func (c *Admin) Get(key string) (admin.Admin, error) {
	var result admin.Admin

	data, err := c.bigCache.Get(c.getPrefixAdmin() + c.getPrefixTelegramID() + key)
	if err != nil {
		return admin.Admin{}, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return admin.Admin{}, err
	}

	return result, nil
}

// Exists checks whether admin exists in BigCache by key.
func (c *Admin) Exists(key string) (bool, error) {
	_, err := c.bigCache.Get(c.getPrefixAdmin() + c.getPrefixTelegramID() + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete removes admin from the cache by key.
func (c *Admin) Delete(key string) error {
	return c.bigCache.Delete(c.getPrefixAdmin() + c.getPrefixTelegramID() + key)
}

// getPrefixAdmin get prefix admin.
func (c *Admin) getPrefixAdmin() string {
	return c.prefixAdmin
}

// getPrefixTelegramID get prefix telegram id.
func (c *Admin) getPrefixTelegramID() string {
	return c.prefixTelegramID
}
