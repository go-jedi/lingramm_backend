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
	GetPrefixTelegramID() string
	Set(key string, val admin.Admin, prefix string) error
	All(prefix string) ([]admin.Admin, error)
	Get(key string, prefix string) (admin.Admin, error)
	Exists(key string, prefix string) (bool, error)
	Delete(key string, prefix string) error
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

// GetPrefixTelegramID get prefix telegram_id.
func (a *Admin) GetPrefixTelegramID() string {
	return a.prefixTelegramID
}

// Set stores admin in BigCache using MessagePack serialization.
func (a *Admin) Set(key string, val admin.Admin, prefix string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return a.bigCache.Set(a.prefixAdmin+prefix+key, b)
}

// All retrieves all admins entries from the cache that match the prefix.
func (a *Admin) All(prefix string) ([]admin.Admin, error) {
	var results []admin.Admin

	iter := a.bigCache.Iterator()
	for iter.SetNext() {
		entry, err := iter.Value()
		if err != nil {
			log.Printf("failed to read cache entry: %v", err)
			continue
		}

		key := entry.Key()

		if !strings.HasPrefix(key, a.prefixAdmin+prefix) {
			continue
		}

		var adm admin.Admin
		if err := msgpack.Unmarshal(entry.Value(), &adm); err != nil {
			log.Printf("failed to unmarshal cache entry for key '%s': %v", key, err)
			continue
		}

		results = append(results, adm)
	}

	return results, nil
}

// Get retrieves admin from BigCache and deserializes it using MessagePack.
func (a *Admin) Get(key string, prefix string) (admin.Admin, error) {
	var result admin.Admin

	data, err := a.bigCache.Get(a.prefixAdmin + prefix + key)
	if err != nil {
		return result, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return admin.Admin{}, err
	}

	return result, nil
}

// Exists checks whether a user exists in BigCache by key.
func (a *Admin) Exists(key string, prefix string) (bool, error) {
	_, err := a.bigCache.Get(a.prefixAdmin + prefix + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete removes admin from the cache by key.
func (a *Admin) Delete(key string, prefix string) error {
	return a.bigCache.Delete(a.prefixAdmin + prefix + key)
}
