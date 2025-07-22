package user

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user/create"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user/exists"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user/exists_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	Create             create.ICreate
	Exists             exists.IExists
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
	GetByTelegramID    getbytelegramid.IGetByTelegramID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create:             create.New(queryTimeout, logger),
		Exists:             exists.New(queryTimeout, logger),
		ExistsByTelegramID: existsbytelegramid.New(queryTimeout, logger),
		GetByTelegramID:    getbytelegramid.New(queryTimeout, logger),
	}
}
