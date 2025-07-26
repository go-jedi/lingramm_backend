package achievementassets

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets/create"
	deletebyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets/delete_by_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	All        all.IAll
	Create     create.ICreate
	DeleteByID deletebyid.IDeleteByID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		All:        all.New(queryTimeout, logger),
		Create:     create.New(queryTimeout, logger),
		DeleteByID: deletebyid.New(queryTimeout, logger),
	}
}
