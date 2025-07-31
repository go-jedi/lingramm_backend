package localizedtext

import (
	createtextcontent "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/create_text_content"
	createtexttranslation "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/create_text_translation"
	deletetextcontentbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/delete_text_content_by_id"
	deletetexttranslationbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/delete_text_translation_by_id"
	existstextcontentbycode "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/exists_text_content_by_code"
	existstextcontentbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/exists_text_content_by_id"
	existstexttranslation "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/exists_text_translation"
	existstexttranslationbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/exists_text_translation_by_id"
	gettextsbylanguage "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text/get_texts_by_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	CreateTextContent         createtextcontent.ICreateTextContent
	CreateTextTranslation     createtexttranslation.ICreateTextTranslation
	DeleteTextContentByID     deletetextcontentbyid.IDeleteTextContentByID
	DeleteTextTranslationByID deletetexttranslationbyid.IDeleteTextTranslationByID
	ExistsTextContentByCode   existstextcontentbycode.IExistsTextContentByCode
	ExistsTextContentByID     existstextcontentbyid.IExistsTextContentByID
	ExistsTextTranslation     existstexttranslation.IExistsTextTranslation
	ExistsTextTranslationByID existstexttranslationbyid.IExistsTextTranslationByID
	GetTextsByLanguage        gettextsbylanguage.IGetTextsByLanguage
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		CreateTextContent:         createtextcontent.New(queryTimeout, logger),
		CreateTextTranslation:     createtexttranslation.New(queryTimeout, logger),
		DeleteTextContentByID:     deletetextcontentbyid.New(queryTimeout, logger),
		DeleteTextTranslationByID: deletetexttranslationbyid.New(queryTimeout, logger),
		ExistsTextContentByCode:   existstextcontentbycode.New(queryTimeout, logger),
		ExistsTextContentByID:     existstextcontentbyid.New(queryTimeout, logger),
		ExistsTextTranslation:     existstexttranslation.New(queryTimeout, logger),
		ExistsTextTranslationByID: existstexttranslationbyid.New(queryTimeout, logger),
		GetTextsByLanguage:        gettextsbylanguage.New(queryTimeout, logger),
	}
}
