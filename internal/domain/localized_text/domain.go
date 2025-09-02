package localizedtext

// TextContents represents text contents in the system.
type TextContents struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Page        string  `json:"page"`
	Description *string `json:"description,omitempty"`
}

// TextTranslations represents text translations in the system.
type TextTranslations struct {
	ID        int64  `json:"id"`
	ContentID int64  `json:"content_id"`
	Lang      string `json:"lang"`
	Value     string `json:"value"`
}

// LocalizedTexts represents localized texts in the system.
type LocalizedTexts struct {
	Code        string  `json:"code"`
	Value       string  `json:"value"`
	Description *string `json:"description,omitempty"`
}

//
// CREATE TEXT CONTENT
//

type CreateTextContentDTO struct {
	Code        string  `json:"code" validate:"required,min=1"`
	Page        string  `json:"page" validate:"required,min=1"`
	Description *string `json:"description" validate:"omitempty,min=1"`
}

//
// CREATE TEXT TRANSLATION
//

type CreateTextTranslationDTO struct {
	ContentID int64  `json:"content_id" validate:"required,gt=0"`
	Lang      string `json:"lang" validate:"required,len=2"`
	Value     string `json:"value" validate:"required,min=1"`
}

//
// SWAGGER
//

type CreateTextContentSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID          int64   `json:"id" example:"1"`
		Code        string  `json:"code" example:"some code"`
		Page        string  `json:"page" example:"some page"`
		Description *string `json:"description,omitempty" example:"some description"`
	} `json:"data"`
}

type CreateTextTranslationSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID        int64  `json:"id" example:"1"`
		ContentID int64  `json:"content_id" example:"1"`
		Lang      string `json:"lang" example:"en"`
		Value     string `json:"value" example:"some value"`
	} `json:"data"`
}

type GetTextsByLanguageSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    map[string][]struct {
		Code        string  `json:"code" example:"some code"`
		Value       string  `json:"value" example:"some value"`
		Description *string `json:"description,omitempty" example:"some description"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
