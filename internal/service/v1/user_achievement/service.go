package userachievement

import (
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	alldetailbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_achievement/all_detail_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	AllDetailByTelegramID alldetailbytelegramid.IAllDetailByTelegramID
}

func New(
	userAchievementRepository *userachievementrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		AllDetailByTelegramID: alldetailbytelegramid.New(userAchievementRepository, userRepository, logger, postgres),
	}
}
