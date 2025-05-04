package create

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	clientassets "github.com/go-jedi/lingvogramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/file_server/client_assets"
	fileserver "github.com/go-jedi/lingvogramm_backend/pkg/file_server"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
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
	fileServer             *fileserver.FileServer
}

func New(
	clientAssetsRepository *clientassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Create {
	return &Create{
		clientAssetsRepository: clientAssetsRepository,
		logger:                 logger,
		postgres:               postgres,
		fileServer:             fileServer,
	}
}

func (c *Create) Execute(ctx context.Context, file *multipart.FileHeader) (clientassets.ClientAssets, error) {
	c.logger.Debug("[create a client assets] execute service")

	var err error

	tx, err := c.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
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

	// convert png or jpg image to webp and upload.
	imageData, err := c.fileServer.ClientAssets.UploadAndConvertToWebP(ctx, file)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	// create client assets.
	result, err := c.clientAssetsRepository.Create.Execute(ctx, tx, imageData)
	if err != nil {
		// compensating action - delete the saved image.
		if err := os.Remove(imageData.ServerPathFile); err != nil {
			c.logger.Warn("failed to remove image after db error", "warn", err)
		}
		return clientassets.ClientAssets{}, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return clientassets.ClientAssets{}, err
	}

	return result, nil
}
