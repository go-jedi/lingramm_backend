package existstextcontentbycode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsTextContentByCode --output=mocks --case=underscore
type IExistsTextContentByCode interface {
	Execute(ctx context.Context, tx pgx.Tx, code string) (bool, error)
}

type ExistsTextContentByCode struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsTextContentByCode {
	r := &ExistsTextContentByCode{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsTextContentByCode) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsTextContentByCode) Execute(ctx context.Context, tx pgx.Tx, code string) (bool, error) {
	r.logger.Debug("[check text content exists by code] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM text_contents
			WHERE code = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		code,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check text content exists by code", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check text content exists by code", "err", err)
		return false, fmt.Errorf("could not check text content exists by code: %w", err)
	}

	return ie, nil
}
