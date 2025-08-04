package deletebyid

import (
	"context"
	"log"
	"os"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteByID --output=mocks --case=underscore
type IDeleteByID interface {
	Execute(ctx context.Context, id int64) (clientassets.ClientAssets, error)
}

type DeleteByID struct {
	clientAssetsRepository *clientassetsrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
	redis                  *redis.Redis
}

func New(
	clientAssetsRepository *clientassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *DeleteByID {
	return &DeleteByID{
		clientAssetsRepository: clientAssetsRepository,
		logger:                 logger,
		postgres:               postgres,
		redis:                  redis,
	}
}

func (s *DeleteByID) Execute(ctx context.Context, id int64) (clientassets.ClientAssets, error) {
	s.logger.Debug("[delete client assets by id] execute service")

	var (
		err    error
		result clientassets.ClientAssets
		ie     bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return clientassets.ClientAssets{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check client assets exists by id.
	ie, err = s.clientAssetsRepository.ExistsByID.Execute(ctx, tx, id)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	if !ie { // if client assets does not exist.
		err = apperrors.ErrClientAssetsDoesNotExist
		return clientassets.ClientAssets{}, err
	}

	// delete client assets by id.
	result, err = s.clientAssetsRepository.DeleteByID.Execute(ctx, tx, id)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	// remove file.
	s.deleteClientFile(ctx, result.NameFileWithoutExtension, result.ServerPathFile)

	err = tx.Commit(ctx)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	return result, nil
}

func (s *DeleteByID) deleteClientFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove client file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileClient.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete client file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("client file removed", "path", path)
	}
}
