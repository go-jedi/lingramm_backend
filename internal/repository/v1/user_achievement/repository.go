package userachievement

import (
	alldetailbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement/all_detail_by_telegram_id"
	unlockavailableachievements "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement/unlock_available_achievements"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AllDetailByTelegramID       alldetailbytelegramid.IAllDetailByTelegramID
	UnlockAvailableAchievements unlockavailableachievements.IUnlockAvailableAchievements
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AllDetailByTelegramID:       alldetailbytelegramid.New(queryTimeout, logger),
		UnlockAvailableAchievements: unlockavailableachievements.New(queryTimeout, logger),
	}
}
