package alldetail

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllDetail --output=mocks --case=underscore
type IAllDetail interface {
	Execute(ctx context.Context) ([]achievement.Detail, error)
}

type AllDetail struct {
	achievementRepository *achievementrepository.Repository
	logger                logger.ILogger
	postgres              *postgres.Postgres
}

func New(
	achievementRepository *achievementrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *AllDetail {
	return &AllDetail{
		achievementRepository: achievementRepository,
		logger:                logger,
		postgres:              postgres,
	}
}

func (s *AllDetail) Execute(ctx context.Context) ([]achievement.Detail, error) {
	s.logger.Debug("[get all detail] execute service")

	var (
		err    error
		result []achievement.Detail
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

	result, err = s.achievementRepository.AllDetail.Execute(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
