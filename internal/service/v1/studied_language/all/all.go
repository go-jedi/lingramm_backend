package all

import (
	"context"
	"log"

	studiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/studied_language"
	studiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context) ([]studiedlanguage.StudiedLanguage, error)
}

type All struct {
	studiedLanguageRepository *studiedlanguagerepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	studiedLanguageRepository *studiedlanguagerepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *All {
	return &All{
		studiedLanguageRepository: studiedLanguageRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *All) Execute(ctx context.Context) ([]studiedlanguage.StudiedLanguage, error) {
	s.logger.Debug("[get all studied languages] execute service")

	var (
		err    error
		result []studiedlanguage.StudiedLanguage
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

	// get all studied languages from database postgres.
	result, err = s.studiedLanguageRepository.All.Execute(ctx, tx)
	if err != nil {
		return nil, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
