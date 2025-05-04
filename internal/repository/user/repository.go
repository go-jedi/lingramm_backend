package user

import (
	"github.com/go-jedi/lingvogramm_backend/internal/repository/user/create"
	"github.com/go-jedi/lingvogramm_backend/internal/repository/user/exists"
	getbytelegramid "github.com/go-jedi/lingvogramm_backend/internal/repository/user/get_by_telegram_id"
	getbyuuid "github.com/go-jedi/lingvogramm_backend/internal/repository/user/get_by_uuid"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
)

type Repository struct {
	Create          create.ICreate
	Exists          exists.IExists
	GetByTelegramID getbytelegramid.IGetByTelegramID
	GetByUUID       getbyuuid.IGetByUUID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create:          create.New(queryTimeout, logger),
		Exists:          exists.New(queryTimeout, logger),
		GetByTelegramID: getbytelegramid.New(queryTimeout, logger),
		GetByUUID:       getbyuuid.New(queryTimeout, logger),
	}
}
