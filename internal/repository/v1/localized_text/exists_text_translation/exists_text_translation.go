package existstexttranslation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsTextTranslation --output=mocks --case=underscore
type IExistsTextTranslation interface {
	Execute(ctx context.Context, tx pgx.Tx, contentID int64, language string) (bool, error)
}

type ExistsTextTranslation struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsTextTranslation {
	r := &ExistsTextTranslation{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsTextTranslation) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsTextTranslation) Execute(ctx context.Context, tx pgx.Tx, contentID int64, language string) (bool, error) {
	r.logger.Debug("[check text translation exists] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM text_translations
			WHERE content_id = $1
			AND lang = $2
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		contentID, language,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check text translation exists", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check text translation exists", "err", err)
		return false, fmt.Errorf("could not check text translation exists: %w", err)
	}

	return ie, nil
}
