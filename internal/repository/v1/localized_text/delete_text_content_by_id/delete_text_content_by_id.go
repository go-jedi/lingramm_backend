package deletetextcontentbyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteTextContentByID --output=mocks --case=underscore
type IDeleteTextContentByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (localizedtext.TextContents, error)
}

type DeleteTextContentByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *DeleteTextContentByID {
	r := &DeleteTextContentByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *DeleteTextContentByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *DeleteTextContentByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (localizedtext.TextContents, error) {
	r.logger.Debug("[delete text content by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM text_contents
		WHERE id = $1
		RETURNING *;
	`

	var dtc localizedtext.TextContents

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(&dtc); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete text content by id", "err", err)
			return localizedtext.TextContents{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete text content by id", "err", err)
		return localizedtext.TextContents{}, fmt.Errorf("could not delete text content by id: %w", err)
	}

	return dtc, nil
}
