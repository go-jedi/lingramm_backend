package clientassets

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/create"
	deletebyid "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/delete_by_id"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	All        all.IAll
	Create     create.ICreate
	DeleteByID deletebyid.IDeleteByID
}

func New(
	clientAssetsRepository *clientassets.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	fileServer *fileserver.FileServer,
) *Service {
	return &Service{
		All:        all.New(clientAssetsRepository, logger, postgres),
		Create:     create.New(clientAssetsRepository, logger, postgres, redis, fileServer),
		DeleteByID: deletebyid.New(clientAssetsRepository, logger, postgres, redis),
	}
}
