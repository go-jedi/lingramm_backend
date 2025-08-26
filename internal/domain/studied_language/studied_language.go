package studiedlanguage

import "time"

type StudiedLanguage struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Lang        string    `json:"lang"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
	Lang        string `json:"lang" validate:"required,len=2"`
}
