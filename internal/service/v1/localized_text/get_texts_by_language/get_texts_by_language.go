package gettextsbylanguage

import (
	"context"
	"fmt"
	"log"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetTextsByLanguage --output=mocks --case=underscore
type IGetTextsByLanguage interface {
	Execute(ctx context.Context, language string) (map[string][]localizedtext.LocalizedTexts, error)
}

type GetTextsByLanguage struct {
	localizedTextRepository *localizedtextepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
	bigCache                *bigcachepkg.BigCache
}

func New(
	localizedTextRepository *localizedtextepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *GetTextsByLanguage {
	return &GetTextsByLanguage{
		localizedTextRepository: localizedTextRepository,
		logger:                  logger,
		postgres:                postgres,
		bigCache:                bigCache,
	}
}

func (s *GetTextsByLanguage) Execute(ctx context.Context, language string) (map[string][]localizedtext.LocalizedTexts, error) {
	s.logger.Debug("[get texts by language] execute service")

	var (
		err    error
		result map[string][]localizedtext.LocalizedTexts
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// get texts by language from cache or database.
	result, err = s.getTextsByLanguageFromCacheOrDatabase(ctx, tx, language)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getTextsByLanguageFromCacheOrDatabase get texts by language from cache or database.
func (s *GetTextsByLanguage) getTextsByLanguageFromCacheOrDatabase(ctx context.Context, tx pgx.Tx, language string) (map[string][]localizedtext.LocalizedTexts, error) {
	// Get localized text by language from cache.
	// If found and no error occurred, return data.
	dataFromCache, err := s.bigCache.LocalizedText.Get(language)
	if err == nil && dataFromCache != nil && len(dataFromCache) > 0 {
		return dataFromCache, nil
	}

	// If localized text by language from db.
	// If found and no error occurred, return data.
	dataFromDB, err := s.localizedTextRepository.GetTextsByLanguage.Execute(ctx, tx, language)
	if err != nil {
		return nil, err
	}

	// set localized text by language in cache.
	if err := s.bigCache.LocalizedText.Set(language, dataFromDB); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to set localized text to cache for language=%s: %v", language, err))
	}

	return dataFromDB, nil
}
