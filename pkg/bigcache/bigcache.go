package bigcache

import (
	"log"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/go-jedi/lingvogramm_backend/pkg/bigcache/iterator"
	"github.com/go-jedi/lingvogramm_backend/pkg/bigcache/user"
)

type BigCache struct {
	User     user.IUser
	Iterator iterator.IIterator
	bigCache *bigcache.BigCache
}

func New(cfg config.BigCacheConfig) (*BigCache, error) {
	bc := &BigCache{}

	bigCacheConfig := bigcache.Config{
		Shards:             cfg.Shards,
		LifeWindow:         time.Duration(cfg.LifeWindow) * time.Minute,
		CleanWindow:        time.Duration(cfg.CleanWindow) * time.Minute,
		MaxEntriesInWindow: cfg.MaxEntriesInWindow,
		MaxEntrySize:       cfg.MaxEntrySize,
		HardMaxCacheSize:   cfg.HardMaxCacheSize,
		Verbose:            cfg.Verbose,
	}
	if cfg.IsOnRemoveWithReason {
		bigCacheConfig.OnRemoveWithReason = func(key string, _ []byte, reason bigcache.RemoveReason) {
			log.Printf("removed key: %s, reason: %v\n", key, reason)
		}
	}

	bigCache, err := bigcache.NewBigCache(bigCacheConfig)
	if err != nil {
		return nil, err
	}

	bc.bigCache = bigCache
	bc.User = user.New(bigCache)
	bc.Iterator = iterator.New(bigCache)

	return bc, nil
}
