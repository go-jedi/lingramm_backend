package gettextsbylanguage

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

//go:generate mockery --name=IGetTextsByLanguage --output=mocks --case=underscore
type IGetTextsByLanguage interface {
	Execute(ctx context.Context, tx pgx.Tx, language string) (map[string][]localizedtext.LocalizedTexts, error)
}

type GetTextsByLanguage struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetTextsByLanguage {
	r := &GetTextsByLanguage{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetTextsByLanguage) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetTextsByLanguage) Execute(ctx context.Context, tx pgx.Tx, language string) (map[string][]localizedtext.LocalizedTexts, error) {
	r.logger.Debug("[get texts by language] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT JSONB_OBJECT_AGG(page, texts) AS result
		FROM (
		    SELECT tc.page, JSONB_AGG(
		    	jsonb_build_object(
		    		'code', tc.code,
		    		'value', tt.value,
            		'description', tc.description
		    	)
		    ) AS texts
		    FROM text_contents tc
		    INNER JOIN text_translations tt ON tc.id = tt.content_id
			WHERE tt.lang = $1
			GROUP BY tc.page
		) sub;
	`

	var lt map[string][]localizedtext.LocalizedTexts

	if err := tx.QueryRow(
		ctxTimeout, q,
		language,
	).Scan(&lt); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get texts by language", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get texts by language", "err", err)
		return nil, fmt.Errorf("could not get texts by language: %w", err)
	}

	return lt, nil
}
