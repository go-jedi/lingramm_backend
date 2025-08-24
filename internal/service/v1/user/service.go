package user

import (
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	GetByTelegramID getbytelegramid.IGetByTelegramID
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		GetByTelegramID: getbytelegramid.New(userRepository, logger, postgres),
	}
}
