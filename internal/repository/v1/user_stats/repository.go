package userstats

import (
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/exists_by_telegram_id"
	getlevelbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/get_level_by_telegram_id"
	getlevelinfobytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/get_level_info_by_telegram_id"
	syncuserstatsfromxpeventsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats/sync_user_stats_from_xp_events_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	ExistsByTelegramID                    existsbytelegramid.IExistsByTelegramID
	GetLevelByTelegramID                  getlevelbytelegramid.IGetLevelByTelegramID
	GetLevelInfoByTelegramID              getlevelinfobytelegramid.IGetLevelInfoByTelegramID
	SyncUserStatsFromXPEventsByTelegramID syncuserstatsfromxpeventsbytelegramid.ISyncUserStatsFromXPEventsByTelegramID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		ExistsByTelegramID:                    existsbytelegramid.New(queryTimeout, logger),
		GetLevelByTelegramID:                  getlevelbytelegramid.New(queryTimeout, logger),
		GetLevelInfoByTelegramID:              getlevelinfobytelegramid.New(queryTimeout, logger),
		SyncUserStatsFromXPEventsByTelegramID: syncuserstatsfromxpeventsbytelegramid.New(queryTimeout, logger),
	}
}
