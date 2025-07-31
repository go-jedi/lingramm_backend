package localizedtext

import (
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	createtextcontent "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text/create_text_content"
	createtexttranslation "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text/create_text_translation"
	gettextsbylanguage "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text/get_texts_by_language"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	CreateTextContent     createtextcontent.ICreateTextContent
	CreateTextTranslation createtexttranslation.ICreateTextTranslation
	GetTextsByLanguage    gettextsbylanguage.IGetTextsByLanguage
}

func New(
	localizedTextRepository *localizedtextepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *Service {
	return &Service{
		CreateTextContent:     createtextcontent.New(localizedTextRepository, logger, postgres),
		CreateTextTranslation: createtexttranslation.New(localizedTextRepository, logger, postgres, bigCache),
		GetTextsByLanguage:    gettextsbylanguage.New(localizedTextRepository, logger, postgres, bigCache),
	}
}
