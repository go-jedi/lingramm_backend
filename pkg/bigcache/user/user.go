package user

import (
	"errors"
	"log"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	"github.com/vmihailenco/msgpack/v5"
)

//go:generate mockery --name=IUser --output=mocks --case=underscore
type IUser interface {
	Set(key string, val user.User) error
	All() ([]user.User, error)
	Get(key string) (user.User, error)
	Exists(key string) (bool, error)
	Delete(key string) error
}

type User struct {
	prefix   string
	bigCache *bigcache.BigCache
}

func New(bigCache *bigcache.BigCache) *User {
	return &User{
		prefix:   "user:",
		bigCache: bigCache,
	}
}

// Set stores a user in BigCache using MessagePack serialization.
func (u *User) Set(key string, val user.User) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return u.bigCache.Set(u.prefix+key, b)
}

// All retrieves all user entries from the cache that match the prefix.
func (u *User) All() ([]user.User, error) {
	var results []user.User

	iter := u.bigCache.Iterator()
	for iter.SetNext() {
		entry, err := iter.Value()
		if err != nil {
			log.Printf("failed to read cache entry: %v", err)
			continue
		}

		key := entry.Key()

		if !strings.HasPrefix(key, u.prefix) {
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
func (u *User) Get(key string) (user.User, error) {
	var result user.User

	data, err := u.bigCache.Get(u.prefix + key)
	if err != nil {
		return result, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

// Exists checks whether a user exists in BigCache by key.
func (u *User) Exists(key string) (bool, error) {
	_, err := u.bigCache.Get(u.prefix + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete removes a user from the cache by key.
func (u *User) Delete(key string) error {
	return u.bigCache.Delete(u.prefix + key)
}
