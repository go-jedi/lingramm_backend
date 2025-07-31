package createtexttranslation

import (
	"context"
	"fmt"
	"log"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateTextTranslation --output=mocks --case=underscore
type ICreateTextTranslation interface {
	Execute(ctx context.Context, dto localizedtext.CreateTextTranslationDTO) (localizedtext.TextTranslations, error)
}

type CreateTextTranslation struct {
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
) *CreateTextTranslation {
	return &CreateTextTranslation{
		localizedTextRepository: localizedTextRepository,
		logger:                  logger,
		postgres:                postgres,
		bigCache:                bigCache,
	}
}

func (s *CreateTextTranslation) Execute(ctx context.Context, dto localizedtext.CreateTextTranslationDTO) (localizedtext.TextTranslations, error) {
	s.logger.Debug("[create text translation] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return localizedtext.TextTranslations{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check text content exists by id.
	existsTextContent, err := s.localizedTextRepository.ExistsTextContentByID.Execute(ctx, tx, dto.ContentID)
	if err != nil {
		return localizedtext.TextTranslations{}, err
	}

	if !existsTextContent { // if text content does not exist.
		return localizedtext.TextTranslations{}, apperrors.ErrTextContentDoesNotExist
	}

	// check text translation exists.
	existsTextTranslation, err := s.localizedTextRepository.ExistsTextTranslation.Execute(ctx, tx, dto.ContentID, dto.Lang)
	if err != nil {
		return localizedtext.TextTranslations{}, err
	}

	if existsTextTranslation { // if text translation already exist.
		return localizedtext.TextTranslations{}, apperrors.ErrTextTranslationAlreadyExists
	}

	// create text translation.
	result, err := s.localizedTextRepository.CreateTextTranslation.Execute(ctx, tx, dto)
	if err != nil {
		return localizedtext.TextTranslations{}, err
	}

	// delete localized text by language in cache.
	if err := s.bigCache.LocalizedText.Delete(dto.Lang); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to delete localized text cache for language=%s: %v", dto.Lang, err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return localizedtext.TextTranslations{}, err
	}

	return result, nil
}
