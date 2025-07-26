package clientassets

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	All    all.IAll
	Create create.ICreate
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		All:    all.New(queryTimeout, logger),
		Create: create.New(queryTimeout, logger),
	}
}
