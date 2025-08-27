package all

import (
	"context"
	"log"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context) ([]eventtype.EventType, error)
}

type All struct {
	eventTypeRepository *eventtyperepository.Repository
	logger              logger.ILogger
	postgres            *postgres.Postgres
}

func New(
	eventTypeRepository *eventtyperepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *All {
	return &All{
		eventTypeRepository: eventTypeRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *All) Execute(ctx context.Context) ([]eventtype.EventType, error) {
	s.logger.Debug("[get all event types] execute service")

	var (
		err    error
		result []eventtype.EventType
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

	// get all event types.
	result, err = s.eventTypeRepository.All.Execute(ctx, tx)
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
