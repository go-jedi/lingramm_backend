package bigcache

import (
	"log"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/pkg/bigcache/admin"
	"github.com/go-jedi/lingramm_backend/pkg/bigcache/iterator"
	localizedtext "github.com/go-jedi/lingramm_backend/pkg/bigcache/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/bigcache/user"
)

type BigCache struct {
	Admin         admin.IAdmin
	Iterator      iterator.IIterator
	LocalizedText localizedtext.ILocalizedText
	User          user.IUser
	bigCache      *bigcache.BigCache
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

	bc.Admin = admin.New(bigCache)
	bc.Iterator = iterator.New(bigCache)
	bc.LocalizedText = localizedtext.New(bigCache)
	bc.User = user.New(bigCache)

	return bc, nil
}
