package studiedlanguage

import (
	studiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language/all"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	All    all.IAll
	Create create.ICreate
}

func New(
	studiedLanguageRepository *studiedlanguagerepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		All:    all.New(studiedLanguageRepository, logger, postgres),
		Create: create.New(studiedLanguageRepository, logger, postgres),
	}
}
