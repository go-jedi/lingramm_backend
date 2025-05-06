package bigcache

import (
	"github.com/go-jedi/lingvogramm_backend/internal/service/bigcache/iterator"
	bigcachepkg "github.com/go-jedi/lingvogramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
)

type Service struct {
	Iterator *iterator.Iterator
}

func New(
	logger logger.ILogger,
	bigCache *bigcachepkg.BigCache,
) *Service {
	return &Service{
		Iterator: iterator.New(logger, bigCache),
	}
}
