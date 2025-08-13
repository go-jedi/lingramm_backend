package subscription

import (
	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
	GetByTelegramID    getbytelegramid.IGetByTelegramID
}

func New(
	subscriptionRepository *subscriptionrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		ExistsByTelegramID: existsbytelegramid.New(subscriptionRepository, userRepository, logger, postgres),
		GetByTelegramID:    getbytelegramid.New(subscriptionRepository, userRepository, logger, postgres),
	}
}
