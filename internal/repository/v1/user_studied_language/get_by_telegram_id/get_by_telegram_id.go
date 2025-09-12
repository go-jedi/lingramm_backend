package getbytelegramid

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

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstudiedlanguage.GetByTelegramIDResponse, error)
}

type GetByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByTelegramID {
	r := &GetByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstudiedlanguage.GetByTelegramIDResponse, error) {
	r.logger.Debug("[get user studied language by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT
			usl.id,
			usl.studied_language_id,
			usl.telegram_id,
			sl.name,
			sl.description,
			sl.lang,
			usl.created_at,
			usl.updated_at
		FROM user_studied_languages usl
		INNER JOIN studied_languages sl ON usl.studied_language_id = sl.id
		WHERE telegram_id = $1;
	`

	var response userstudiedlanguage.GetByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(
		&response.ID, &response.StudiedLanguageID,
		&response.TelegramID, &response.Name,
		&response.Description, &response.Lang,
		&response.CreatedAt, &response.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get user studied language by telegram id", "err", err)
			return userstudiedlanguage.GetByTelegramIDResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get user studied language by telegram id", "err", err)
		return userstudiedlanguage.GetByTelegramIDResponse{}, fmt.Errorf("could not get user studied language by telegram id: %w", err)
	}

	return response, nil
}
