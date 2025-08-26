package create

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

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto studiedlanguage.CreateDTO) (studiedlanguage.StudiedLanguage, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	r := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Create) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto studiedlanguage.CreateDTO) (studiedlanguage.StudiedLanguage, error) {
	r.logger.Debug("[create a new studied language] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO studied_languages(
		    name,
		    description,
		    lang
		) VALUES ($1, $2, $3)
		RETURNING *;
	`

	var nsl studiedlanguage.StudiedLanguage

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.Name, dto.Description, dto.Lang,
	).Scan(
		&nsl.ID, &nsl.Name, &nsl.Description,
		&nsl.Lang, &nsl.CreatedAt, &nsl.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new studied language", "err", err)
			return studiedlanguage.StudiedLanguage{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new studied language", "err", err)
		return studiedlanguage.StudiedLanguage{}, fmt.Errorf("could not create a new studied language: %w", err)
	}

	return nsl, nil
}
