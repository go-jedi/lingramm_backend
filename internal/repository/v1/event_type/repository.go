package eventtype

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type/all"
	allbynames "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type/all_by_names"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type/create"
	existsbyname "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type/exists_by_name"
	getbyname "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type/get_by_name"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	All          all.IAll
	AllByNames   allbynames.IAllByNames
	Create       create.ICreate
	ExistsByName existsbyname.IExistsByName
	GetByName    getbyname.IGetByName
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		All:          all.New(queryTimeout, logger),
		AllByNames:   allbynames.New(queryTimeout, logger),
		Create:       create.New(queryTimeout, logger),
		ExistsByName: existsbyname.New(queryTimeout, logger),
		GetByName:    getbyname.New(queryTimeout, logger),
	}
}
