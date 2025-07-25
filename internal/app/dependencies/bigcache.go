package dependencies

import (
	bigcachehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache"
	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
)

func (d *Dependencies) BigCacheService() *bigcacheservice.Service {
	if d.bigCacheService == nil {
		d.bigCacheService = bigcacheservice.New(
			d.logger,
			d.bigCache,
		)
	}

	return d.bigCacheService
}

func (d *Dependencies) BigCacheHandler() *bigcachehandler.Handler {
	if d.bigCacheHandler == nil {
		d.bigCacheHandler = bigcachehandler.New(
			d.BigCacheService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.bigCacheHandler
}
