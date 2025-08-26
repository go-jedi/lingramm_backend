package update

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

//go:generate mockery --name=IUpdate --output=mocks --case=underscore
type IUpdate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto userstudiedlanguage.UpdateDTO) (userstudiedlanguage.UserStudiedLanguage, error)
}

type Update struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Update {
	r := &Update{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Update) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Update) Execute(ctx context.Context, tx pgx.Tx, dto userstudiedlanguage.UpdateDTO) (userstudiedlanguage.UserStudiedLanguage, error) {
	r.logger.Debug("[update user studied language] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		UPDATE user_studied_languages SET
		    studied_languages_id = $1
		WHERE telegram_id = $2;
	`

	var usl userstudiedlanguage.UserStudiedLanguage

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.StudiedLanguagesID, dto.TelegramID,
	).Scan(
		&usl.ID, &usl.TelegramID, &usl.StudiedLanguagesID,
		&usl.CreatedAt, &usl.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while update user studied language", "err", err)
			return userstudiedlanguage.UserStudiedLanguage{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to update user studied language", "err", err)
		return userstudiedlanguage.UserStudiedLanguage{}, fmt.Errorf("could not update user studied language: %w", err)
	}

	return usl, nil
}
