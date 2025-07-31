package createtexttranslation

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

//go:generate mockery --name=ICreateTextTranslation --output=mocks --case=underscore
type ICreateTextTranslation interface {
	Execute(ctx context.Context, tx pgx.Tx, dto localizedtext.CreateTextTranslationDTO) (localizedtext.TextTranslations, error)
}

type CreateTextTranslation struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateTextTranslation {
	r := &CreateTextTranslation{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateTextTranslation) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateTextTranslation) Execute(ctx context.Context, tx pgx.Tx, dto localizedtext.CreateTextTranslationDTO) (localizedtext.TextTranslations, error) {
	r.logger.Debug("[create text translation] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO text_translations(
		    content_id,
		    lang,
		    value
		) VALUES($1, $2, $3)
		RETURNING *;
	`

	var ntt localizedtext.TextTranslations

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.ContentID, dto.Lang,
		dto.Value,
	).Scan(
		&ntt.ID, &ntt.ContentID,
		&ntt.Lang, &ntt.Value,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create text translation", "err", err)
			return localizedtext.TextTranslations{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create text translation", "err", err)
		return localizedtext.TextTranslations{}, fmt.Errorf("could not create text translation: %w", err)
	}

	return ntt, nil
}
