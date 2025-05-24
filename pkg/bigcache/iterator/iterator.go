package iterator

import (
	"context"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	"github.com/vmihailenco/msgpack/v5"
)

// TypeFactory is a function that returns a new empty instance of a specific type.
// Used to dynamically create values for unmarshaling from cache.
type TypeFactory func() any

//go:generate mockery --name=IIterator --output=mocks --case=underscore
type IIterator interface {
	RegisterType(prefix string, factory TypeFactory)
	Iterator(ctx context.Context) (map[string]any, error)
}

// Iterator implements IIterator and provides functionality
// to iterate and decode values from BigCache using registered type factories.
type Iterator struct {
	bigCache   *bigcache.BigCache
	typeLookup map[string]TypeFactory
}

// New creates a new Iterator with predefined type factories.
func New(bigCache *bigcache.BigCache) *Iterator {
	return &Iterator{
		bigCache: bigCache,
		typeLookup: map[string]TypeFactory{
			"user:":  func() any { return new(user.User) },
			"admin:": func() any { return new(admin.Admin) },
		},
	}
}

// RegisterType allows dynamic registration of a new type factory
// to handle decoding cache entries with a specific key prefix.
func (i *Iterator) RegisterType(prefix string, factory TypeFactory) {
	i.typeLookup[prefix] = factory
}

// Iterator retrieves all entries from BigCache, detects their type using prefix,
// and decodes them using MessagePack. It returns a map of decoded objects.
// Context cancellation is respected during iteration.
func (i *Iterator) Iterator(ctx context.Context) (map[string]any, error) {
	// create an iterator over the cache.
	iterator := i.bigCache.Iterator()
	entries := make(map[string]any)

	// loop through all cache entries.
	for iterator.SetNext() {
		// check if the context has been canceled or exceeded deadline.
		if err := ctx.Err(); err != nil {
			return nil, err // context.Canceled or context.DeadlineExceeded.
		}

		// retrieve the current cache entry.
		entry, err := iterator.Value()
		if err != nil {
			return nil, err
		}

		key := entry.Key()

		// detect appropriate factory based on key prefix.
		var value any
		for k := range i.typeLookup {
			if strings.HasPrefix(key, k) {
				value = i.typeLookup[k]()
				break
			}
		}

		// skip unknown types (no registered factory for prefix).
		if value == nil {
			continue // type is not register.
		}

		// decode the MessagePack-encoded value into the typed object.
		if err := msgpack.Unmarshal(entry.Value(), value); err != nil {
			return nil, err
		}

		entries[key] = value
	}

	return entries, nil
}
