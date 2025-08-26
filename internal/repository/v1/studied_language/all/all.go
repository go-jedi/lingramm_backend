package all

import (
	"context"
	"errors"
	"fmt"
	"time"

	studiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]studiedlanguage.StudiedLanguage, error)
}

type All struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *All {
	r := &All{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *All) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *All) Execute(ctx context.Context, tx pgx.Tx) ([]studiedlanguage.StudiedLanguage, error) {
	r.logger.Debug("[get all studied languages] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT * 
		FROM studied_languages
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all studied languages", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all studied languages", "err", err)
		return nil, fmt.Errorf("could not get all studied languages: %w", err)
	}
	defer rows.Close()

	var studiedLanguages []studiedlanguage.StudiedLanguage

	for rows.Next() {
		var studiedLanguage studiedlanguage.StudiedLanguage

		if err := rows.Scan(
			&studiedLanguage.ID, &studiedLanguage.Name,
			&studiedLanguage.Description, &studiedLanguage.Lang,
			&studiedLanguage.CreatedAt, &studiedLanguage.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all studied languages", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all studied languages: %w", err)
		}

		studiedLanguages = append(studiedLanguages, studiedLanguage)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all studied languages", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all studied languages: %w", err)
	}

	return studiedLanguages, nil
}
