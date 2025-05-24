package user

import (
	"errors"
	"log"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	prefixUser       = "user:"
	prefixTelegramID = "telegram_id:"
	prefixUUID       = "uuid:"
)

//go:generate mockery --name=IUser --output=mocks --case=underscore
type IUser interface {
	GetPrefixTelegramID() string
	GetPrefixUUID() string
	Set(key string, val user.User, prefix string) error
	All(prefix string) ([]user.User, error)
	Get(key string, prefix string) (user.User, error)
	Exists(key string, prefix string) (bool, error)
	Delete(key string, prefix string) error
}

type User struct {
	prefixUser       string
	prefixTelegramID string
	prefixUUID       string
	bigCache         *bigcache.BigCache
}

func New(bigCache *bigcache.BigCache) *User {
	return &User{
		prefixUser:       prefixUser,
		prefixTelegramID: prefixTelegramID,
		prefixUUID:       prefixUUID,
		bigCache:         bigCache,
	}
}

// GetPrefixTelegramID get prefix telegram_id.
func (u *User) GetPrefixTelegramID() string {
	return u.prefixTelegramID
}

// GetPrefixUUID get prefix uuid.
func (u *User) GetPrefixUUID() string {
	return u.prefixUUID
}

// Set stores a user in BigCache using MessagePack serialization.
func (u *User) Set(key string, val user.User, prefix string) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return u.bigCache.Set(u.prefixUser+prefix+key, b)
}

// All retrieves all user entries from the cache that match the prefix.
func (u *User) All(prefix string) ([]user.User, error) {
	var results []user.User

	iter := u.bigCache.Iterator()
	for iter.SetNext() {
		entry, err := iter.Value()
		if err != nil {
			log.Printf("failed to read cache entry: %v", err)
			continue
		}

		key := entry.Key()

		if !strings.HasPrefix(key, u.prefixUser+prefix) {
			continue
		}

		var usr user.User
		if err := msgpack.Unmarshal(entry.Value(), &usr); err != nil {
			log.Printf("failed to unmarshal cache entry for key '%s': %v", key, err)
			continue
		}

		results = append(results, usr)
	}

	return results, nil
}

// Get retrieves a user from BigCache and deserializes it using MessagePack.
func (u *User) Get(key string, prefix string) (user.User, error) {
	var result user.User

	data, err := u.bigCache.Get(u.prefixUser + prefix + key)
	if err != nil {
		return result, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

// Exists checks whether a user exists in BigCache by key.
func (u *User) Exists(key string, prefix string) (bool, error) {
	_, err := u.bigCache.Get(u.prefixUser + prefix + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete removes a user from the cache by key.
func (u *User) Delete(key string, prefix string) error {
	return u.bigCache.Delete(u.prefixUser + prefix + key)
}
