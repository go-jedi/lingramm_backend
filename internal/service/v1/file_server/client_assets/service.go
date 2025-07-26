package clientassets

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/create"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	All    all.IAll
	Create create.ICreate
}

func New(
	clientAssetsRepository *clientassets.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Service {
	return &Service{
		All:    all.New(clientAssetsRepository, logger, postgres),
		Create: create.New(clientAssetsRepository, logger, postgres, fileServer),
	}
}
