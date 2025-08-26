package userstudiedlanguage

import (
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstudiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language/create"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language/update"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	Create             create.ICreate
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
	GetByTelegramID    getbytelegramid.IGetByTelegramID
	Update             update.IUpdate
}

func New(
	userStudiedLanguageRepository *userstudiedlanguagerepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		Create:             create.New(userStudiedLanguageRepository, userRepository, logger, postgres),
		ExistsByTelegramID: existsbytelegramid.New(userStudiedLanguageRepository, userRepository, logger, postgres),
		GetByTelegramID:    getbytelegramid.New(userStudiedLanguageRepository, userRepository, logger, postgres),
		Update:             update.New(userStudiedLanguageRepository, userRepository, logger, postgres),
	}
}
