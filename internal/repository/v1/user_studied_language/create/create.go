package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	userstudiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto userstudiedlanguage.CreateDTO) (userstudiedlanguage.UserStudiedLanguage, error)
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

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto userstudiedlanguage.CreateDTO) (userstudiedlanguage.UserStudiedLanguage, error) {
	r.logger.Debug("[create a new user studied language] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO user_studied_languages(
		    telegram_id,
		    studied_language_id
		) VALUES($1, $2)
		RETURNING *;
	`

	var nusl userstudiedlanguage.UserStudiedLanguage

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.TelegramID, dto.StudiedLanguageID,
	).Scan(
		&nusl.ID, &nusl.TelegramID, &nusl.StudiedLanguageID,
		&nusl.CreatedAt, &nusl.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new user studied language", "err", err)
			return userstudiedlanguage.UserStudiedLanguage{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new user studied language", "err", err)
		return userstudiedlanguage.UserStudiedLanguage{}, fmt.Errorf("could not create a new user studied language: %w", err)
	}

	return nusl, nil
}
