package studiedlanguage

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language/all"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language/create"
	existsbylang "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language/exists_by_lang"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	All          all.IAll
	Create       create.ICreate
	ExistsByLang existsbylang.IExistsByLang
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		All:          all.New(queryTimeout, logger),
		Create:       create.New(queryTimeout, logger),
		ExistsByLang: existsbylang.New(queryTimeout, logger),
	}
}
