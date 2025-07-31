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
	prefixUser       string
	prefixTelegramID string
	bigCache         *bigcache.BigCache
}

func New(bigCache *bigcache.BigCache) *User {
	return &User{
		prefixUser:       prefixUser,
		prefixTelegramID: prefixTelegramID,
		bigCache:         bigCache,
	}
}

// Set stores a user in BigCache using MessagePack serialization.
func (c *User) Set(key string, val user.User) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return c.bigCache.Set(c.getPrefixUser()+c.getPrefixTelegramID()+key, b)
}

// All retrieves all user entries from the cache that match the prefix.
func (c *User) All() ([]user.User, error) {
	var result []user.User

	iter := c.bigCache.Iterator()
	for iter.SetNext() {
		entry, err := iter.Value()
		if err != nil {
			log.Printf("failed to read cache entry: %v", err)
			continue
		}

		key := entry.Key()

		if !strings.HasPrefix(key, c.getPrefixUser()+c.getPrefixTelegramID()) {
			continue
		}

		var usr user.User
		if err := msgpack.Unmarshal(entry.Value(), &usr); err != nil {
			log.Printf("failed to unmarshal cache entry for key '%s': %v", key, err)
			continue
		}

		result = append(result, usr)
	}

	return result, nil
}

// Get retrieves a user from BigCache and deserializes it using MessagePack.
func (c *User) Get(key string) (user.User, error) {
	var result user.User

	data, err := c.bigCache.Get(c.getPrefixUser() + c.getPrefixTelegramID() + key)
	if err != nil {
		return user.User{}, err
	}

	if err := msgpack.Unmarshal(data, &result); err != nil {
		return user.User{}, err
	}

	return result, nil
}

// Exists checks whether a user exists in BigCache by key.
func (c *User) Exists(key string) (bool, error) {
	_, err := c.bigCache.Get(c.getPrefixUser() + c.getPrefixTelegramID() + key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete removes a user from the cache by key.
func (c *User) Delete(key string) error {
	return c.bigCache.Delete(c.getPrefixUser() + c.getPrefixTelegramID() + key)
}

// getPrefixUser get prefix user.
func (c *User) getPrefixUser() string {
	return c.prefixUser
}

// getPrefixTelegramID get prefix telegram id.
func (c *User) getPrefixTelegramID() string {
	return c.prefixTelegramID
}
