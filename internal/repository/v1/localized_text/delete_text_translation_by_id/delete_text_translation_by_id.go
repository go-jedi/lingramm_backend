package deletetexttranslationbyid

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

//go:generate mockery --name=IDeleteTextTranslationByID --output=mocks --case=underscore
type IDeleteTextTranslationByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (localizedtext.TextTranslations, error)
}

type DeleteTextTranslationByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *DeleteTextTranslationByID {
	r := &DeleteTextTranslationByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *DeleteTextTranslationByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *DeleteTextTranslationByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (localizedtext.TextTranslations, error) {
	r.logger.Debug("[delete text translation by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM text_translations
		WHERE id = $1
		RETURNING *;
	`

	var dtt localizedtext.TextTranslations

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(
		&dtt.ID, &dtt.ContentID,
		&dtt.Lang, &dtt.Value,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete text translation by id", "err", err)
			return localizedtext.TextTranslations{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete text translation by id", "err", err)
		return localizedtext.TextTranslations{}, fmt.Errorf("could not delete text translation by id: %w", err)
	}

	return dtt, nil
}
