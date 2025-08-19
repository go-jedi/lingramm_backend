package userstats

import (
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	getlevelbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats/get_level_by_telegram_id"
	getlevelinfobytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats/get_level_info_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	GetLevelByTelegramID     getlevelbytelegramid.IGetLevelByTelegramID
	GetLevelInfoByTelegramID getlevelinfobytelegramid.IGetLevelInfoByTelegramID
}

func New(
	userStatsRepository *userstatsrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		GetLevelByTelegramID:     getlevelbytelegramid.New(userStatsRepository, userRepository, logger, postgres),
		GetLevelInfoByTelegramID: getlevelinfobytelegramid.New(userStatsRepository, userRepository, logger, postgres),
	}
}
