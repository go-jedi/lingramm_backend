package createtextcontent

import (
	"context"
	"log"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateTextContent --output=mocks --case=underscore
type ICreateTextContent interface {
	Execute(ctx context.Context, dto localizedtext.CreateTextContentDTO) (localizedtext.TextContents, error)
}

type CreateTextContent struct {
	localizedTextRepository *localizedtextepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
}

func New(
	localizedTextRepository *localizedtextepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *CreateTextContent {
	return &CreateTextContent{
		localizedTextRepository: localizedTextRepository,
		logger:                  logger,
		postgres:                postgres,
	}
}

func (s *CreateTextContent) Execute(ctx context.Context, dto localizedtext.CreateTextContentDTO) (localizedtext.TextContents, error) {
	s.logger.Debug("[create text content] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return localizedtext.TextContents{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check text content exists by code.
	ie, err := s.localizedTextRepository.ExistsTextContentByCode.Execute(ctx, tx, dto.Code)
	if err != nil {
		return localizedtext.TextContents{}, err
	}

	if ie { // if text content already exist.
		return localizedtext.TextContents{}, apperrors.ErrTextContentAlreadyExists
	}

	// create text content.
	result, err := s.localizedTextRepository.CreateTextContent.Execute(ctx, tx, dto)
	if err != nil {
		return localizedtext.TextContents{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return localizedtext.TextContents{}, err
	}

	return result, nil
}
