package userstats

import (
	ensurestreakdaysincrementtoday "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/ensure_streak_days_increment_today"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/exists_by_telegram_id"
	getlevelbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/get_level_by_telegram_id"
	getlevelinfobytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/get_level_info_by_telegram_id"
	hasstreakdaysincrementtoday "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/has_streak_days_increment_today"
	syncuserstatsfromxpeventsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/sync_user_stats_from_xp_events_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	EnsureStreakDaysIncrementToday        ensurestreakdaysincrementtoday.IEnsureStreakDaysIncrementToday
	ExistsByTelegramID                    existsbytelegramid.IExistsByTelegramID
	GetLevelByTelegramID                  getlevelbytelegramid.IGetLevelByTelegramID
	GetLevelInfoByTelegramID              getlevelinfobytelegramid.IGetLevelInfoByTelegramID
	HasStreakDaysIncrementToday           hasstreakdaysincrementtoday.IHasStreakDaysIncrementToday
	SyncUserStatsFromXPEventsByTelegramID syncuserstatsfromxpeventsbytelegramid.ISyncUserStatsFromXPEventsByTelegramID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		EnsureStreakDaysIncrementToday:        ensurestreakdaysincrementtoday.New(queryTimeout, logger),
		ExistsByTelegramID:                    existsbytelegramid.New(queryTimeout, logger),
		GetLevelByTelegramID:                  getlevelbytelegramid.New(queryTimeout, logger),
		GetLevelInfoByTelegramID:              getlevelinfobytelegramid.New(queryTimeout, logger),
		HasStreakDaysIncrementToday:           hasstreakdaysincrementtoday.New(queryTimeout, logger),
		SyncUserStatsFromXPEventsByTelegramID: syncuserstatsfromxpeventsbytelegramid.New(queryTimeout, logger),
	}
}
