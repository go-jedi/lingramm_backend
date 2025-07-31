package existstexttranslationbyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsTextTranslationByID --output=mocks --case=underscore
type IExistsTextTranslationByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error)
}

type ExistsTextTranslationByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsTextTranslationByID {
	r := &ExistsTextTranslationByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsTextTranslationByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsTextTranslationByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error) {
	r.logger.Debug("[check text translation exists by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM text_translations
			WHERE id = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check text translation exists by id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check text translation exists by id", "err", err)
		return false, fmt.Errorf("could not check text translation exists by id: %w", err)
	}

	return ie, nil
}
