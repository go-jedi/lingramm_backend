package clientassets

import (
	clientassets "github.com/go-jedi/lingvogramm_backend/internal/repository/file_server/client_assets"
	"github.com/go-jedi/lingvogramm_backend/internal/service/file_server/client_assets/all"
	"github.com/go-jedi/lingvogramm_backend/internal/service/file_server/client_assets/create"
	fileserver "github.com/go-jedi/lingvogramm_backend/pkg/file_server"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
)

type Service struct {
	Create create.ICreate
	All    all.IAll
}

func New(
	clientAssetsRepository *clientassets.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Service {
	return &Service{
		Create: create.New(clientAssetsRepository, logger, postgres, fileServer),
		All:    all.New(clientAssetsRepository, logger, postgres),
	}
}
