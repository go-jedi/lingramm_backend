package getbyname

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

//go:generate mockery --name=IGetByName --output=mocks --case=underscore
type IGetByName interface {
	Execute(ctx context.Context, name string) (eventtype.EventType, error)
}

type GetByName struct {
	eventTypeRepository *eventtyperepository.Repository
	logger              logger.ILogger
	postgres            *postgres.Postgres
}

func New(
	eventTypeRepository *eventtyperepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetByName {
	return &GetByName{
		eventTypeRepository: eventTypeRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *GetByName) Execute(ctx context.Context, name string) (eventtype.EventType, error) {
	s.logger.Debug("[get event type by name] execute service")

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
	ie, err = s.eventTypeRepository.ExistsByName.Execute(ctx, tx, name)
	if err != nil {
		return eventtype.EventType{}, err
	}

	if !ie { // if event type does not exist.
		err = apperrors.ErrEventTypeDoesNotExist
		return eventtype.EventType{}, err
	}

	// get event type by name.
	result, err = s.eventTypeRepository.GetByName.Execute(ctx, tx, name)
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
