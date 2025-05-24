package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, data clientassets.UploadAndConvertToWebpResponse) (clientassets.ClientAssets, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	c := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	c.init()

	return c
}

func (c *Create) init() {
	if c.queryTimeout == 0 {
		c.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (c *Create) Execute(ctx context.Context, tx pgx.Tx, data clientassets.UploadAndConvertToWebpResponse) (clientassets.ClientAssets, error) {
	c.logger.Debug("[create a client assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO client_assets(
			name_file,
		    server_path_file,
		    client_path_file,
		    extension,
		    quality,
		    old_name_file,
		    old_extension
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING *;
	`

	var ca clientassets.ClientAssets

	if err := tx.QueryRow(
		ctxTimeout, q,
		data.NameFile, data.ServerPathFile, data.ClientPathFile,
		data.Extension, data.Quality, data.OldNameFile, data.OldExtension,
	).Scan(
		&ca.ID, &ca.NameFile, &ca.ServerPathFile,
		&ca.ClientPathFile, &ca.Extension, &ca.Quality,
		&ca.OldNameFile, &ca.OldExtension, &ca.CreatedAt, &ca.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.logger.Error("request timed out while create a client assets", "err", err)
			return clientassets.ClientAssets{}, fmt.Errorf("the request timed out: %w", err)
		}
		c.logger.Error("failed to create a client assets", "err", err)
		return clientassets.ClientAssets{}, fmt.Errorf("could not create a client assets: %w", err)
	}

	return ca, nil
}
