package iterator

import (
	"context"

	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

//go:generate mockery --name=IIterator --output=mocks --case=underscore
type IIterator interface {
	Execute(ctx context.Context) (map[string]any, error)
}

type Iterator struct {
	logger   logger.ILogger
	bigCache *bigcachepkg.BigCache
}

func New(
	logger logger.ILogger,
	bigCache *bigcachepkg.BigCache,
) *Iterator {
	return &Iterator{
		logger:   logger,
		bigCache: bigCache,
	}
}

func (i *Iterator) Execute(ctx context.Context) (map[string]any, error) {
	i.logger.Debug("[iterator for show data in bigcache] execute service")

	result, err := i.bigCache.Iterator.Iterator(ctx)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
