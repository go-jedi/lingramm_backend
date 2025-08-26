package create

import (
	"context"
	"log"

	studiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/studied_language"
	studiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto studiedlanguage.CreateDTO) (studiedlanguage.StudiedLanguage, error)
}

type Create struct {
	studiedLanguageRepository *studiedlanguagerepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	studiedLanguageRepository *studiedlanguagerepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Create {
	return &Create{
		studiedLanguageRepository: studiedLanguageRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *Create) Execute(ctx context.Context, dto studiedlanguage.CreateDTO) (studiedlanguage.StudiedLanguage, error) {
	s.logger.Debug("[create a new studied language] execute service")

	var (
		err    error
		result studiedlanguage.StudiedLanguage
		ie     bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return studiedlanguage.StudiedLanguage{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check studied language exists by lang.
	ie, err = s.studiedLanguageRepository.ExistsByLang.Execute(ctx, tx, dto.Lang)
	if err != nil {
		return studiedlanguage.StudiedLanguage{}, err
	}

	if ie { // if studied language already exist.
		err = apperrors.ErrStudiedLanguageAlreadyExists
		return studiedlanguage.StudiedLanguage{}, err
	}

	// create studied language.
	result, err = s.studiedLanguageRepository.Create.Execute(ctx, tx, dto)
	if err != nil {
		return studiedlanguage.StudiedLanguage{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return studiedlanguage.StudiedLanguage{}, err
	}

	return result, nil
}
