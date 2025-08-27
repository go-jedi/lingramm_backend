package create

import (
	"context"
	"log"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto eventtype.CreateDTO) (eventtype.EventType, error)
}

type Create struct {
	eventTypeRepository *eventtyperepository.Repository
	logger              logger.ILogger
	postgres            *postgres.Postgres
}

func New(
	eventTypeRepository *eventtyperepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Create {
	return &Create{
		eventTypeRepository: eventTypeRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *Create) Execute(ctx context.Context, dto eventtype.CreateDTO) (eventtype.EventType, error) {
	s.logger.Debug("[create a new event type] execute service")

	var (
		err    error
		result eventtype.EventType
		ie     bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return eventtype.EventType{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check event type exists by name
	ie, err = s.eventTypeRepository.ExistsByName.Execute(ctx, tx, dto.Name)
	if err != nil {
		return eventtype.EventType{}, err
	}

	if ie { // if event type already exist.
		err = apperrors.ErrEventTypeAlreadyExists
		return eventtype.EventType{}, err
	}

	// create new event type.
	result, err = s.eventTypeRepository.Create.Execute(ctx, tx, dto)
	if err != nil {
		return eventtype.EventType{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return eventtype.EventType{}, err
	}

	return result, nil
}
