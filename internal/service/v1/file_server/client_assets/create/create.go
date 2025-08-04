package create

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, file *multipart.FileHeader) (clientassets.ClientAssets, error)
}

type Create struct {
	clientAssetsRepository *clientassetsrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
	redis                  *redis.Redis
	fileServer             *fileserver.FileServer
}

func New(
	clientAssetsRepository *clientassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	fileServer *fileserver.FileServer,
) *Create {
	return &Create{
		clientAssetsRepository: clientAssetsRepository,
		logger:                 logger,
		postgres:               postgres,
		redis:                  redis,
		fileServer:             fileServer,
	}
}

func (s *Create) Execute(ctx context.Context, file *multipart.FileHeader) (clientassets.ClientAssets, error) {
	s.logger.Debug("[create a client assets] execute service")

	var (
		err       error
		imageData clientassets.UploadAndConvertToWebpResponse
		result    clientassets.ClientAssets
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
			s.deleteClientTempFile(ctx, imageData.NameFileWithoutExtension, imageData.ServerPathFile)
		}
	}()

	// convert png or jpg image to webp and upload.
	imageData, err = s.fileServer.ClientAssets.UploadAndConvertToWebP(ctx, file)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	// create client assets.
	result, err = s.clientAssetsRepository.Create.Execute(ctx, tx, imageData)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	return result, nil
}

// deleteClientTempFile delete client temp file.
func (s *Create) deleteClientTempFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if path == "" {
		return
	}

	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove temporary client file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileClient.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete client file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("successfully removed temporary client file", "path", path)
	}
}
