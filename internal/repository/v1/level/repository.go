package level

import (
	backfillmissinglevelhistorybytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/level/back_fill_missing_level_history_by_telegram_id"
	createuserlevelhistory "github.com/go-jedi/lingramm_backend/internal/repository/v1/level/create_user_level_history"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	BackFillMissingLevelHistoryByTelegramID backfillmissinglevelhistorybytelegramid.IBackFillMissingLevelHistoryByTelegramID
	CreateUserLevelHistory                  createuserlevelhistory.ICreateUserLevelHistory
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		BackFillMissingLevelHistoryByTelegramID: backfillmissinglevelhistorybytelegramid.New(queryTimeout, logger),
		CreateUserLevelHistory:                  createuserlevelhistory.New(queryTimeout, logger),
	}
}
