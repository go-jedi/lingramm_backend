package updatestatus

import (
	"context"
	"log"

	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IUpdateStatus --output=mocks --case=underscore
type IUpdateStatus interface {
	Execute(ctx context.Context, id int64, status string) error
}

type UpdateStatus struct {
	notificationRepository *notificationrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
}

func New(
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *UpdateStatus {
	return &UpdateStatus{
		notificationRepository: notificationRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (s *UpdateStatus) Execute(ctx context.Context, id int64, status string) error {
	s.logger.Debug("[update status notification] execute service")

	var (
		err error
		ie  bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check notification exists by id.
	ie, err = s.notificationRepository.ExistsByID.Execute(ctx, tx, id)
	if err != nil {
		return err
	}

	if !ie { // if notification does not exist.
		err = apperrors.ErrNotificationDoesNotExist
		return err
	}

	// update notification status.
	err = s.notificationRepository.UpdateStatus.Execute(ctx, tx, id, status)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
