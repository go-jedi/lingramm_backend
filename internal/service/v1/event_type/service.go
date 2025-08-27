package eventtype

import (
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/event_type/all"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/event_type/create"
	getbyname "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type/get_by_name"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	All       all.IAll
	Create    create.ICreate
	GetByName getbyname.IGetByName
}

func New(
	eventTypeRepository *eventtyperepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		All:       all.New(eventTypeRepository, logger, postgres),
		Create:    create.New(eventTypeRepository, logger, postgres),
		GetByName: getbyname.New(eventTypeRepository, logger, postgres),
	}
}
