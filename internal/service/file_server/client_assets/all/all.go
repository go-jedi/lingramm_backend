package all

import (
	"context"
	"log"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context) ([]clientassets.ClientAssets, error)
}

type All struct {
	clientAssetsRepository *clientassetsrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
}

func New(
	clientAssetsRepository *clientassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *All {
	return &All{
		clientAssetsRepository: clientAssetsRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (a *All) Execute(ctx context.Context) ([]clientassets.ClientAssets, error) {
	a.logger.Debug("[get all client assets] execute service")

	var err error

	tx, err := a.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
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

	// get all client assets from database postgres.
	result, err := a.clientAssetsRepository.All.Execute(ctx, tx)
	if err != nil {
		return nil, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
