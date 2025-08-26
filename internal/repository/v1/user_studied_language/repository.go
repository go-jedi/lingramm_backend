package userstudiedlanguage

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language/create"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language/update"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	Create             create.ICreate
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
	GetByTelegramID    getbytelegramid.IGetByTelegramID
	Update             update.IUpdate
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create:             create.New(queryTimeout, logger),
		ExistsByTelegramID: existsbytelegramid.New(queryTimeout, logger),
		GetByTelegramID:    getbytelegramid.New(queryTimeout, logger),
		Update:             update.New(queryTimeout, logger),
	}
}
