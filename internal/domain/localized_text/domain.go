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
