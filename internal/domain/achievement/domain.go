package achievement

import (
	"mime/multipart"
	"time"

	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	awardassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/award_assets"
)

// Achievement represents achievement in the system.
type Achievement struct {
	ID                  int64     `json:"id"`
	AchievementAssetsID int64     `json:"achievement_assets_id"`
	AwardAssetsID       int64     `json:"award_assets_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Description         *string   `json:"description,omitempty"`
	Code                string    `json:"code"`
	Name                string    `json:"name"`
}

// Condition represents achievement condition in the system.
type Condition struct {
	ID            int64     `json:"id"`
	AchievementID int64     `json:"achievement_id"`
	Value         int64     `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ConditionType string    `json:"condition_type"`
	Operator      string    `json:"operator"`
}

// Detail represents achievement detail in the system.
type Detail struct {
	Achievement       Achievement                         `json:"achievement"`
	Condition         Condition                           `json:"condition"`
	AchievementAssets achievementassets.AchievementAssets `json:"achievement_assets"`
	AwardAssets       awardassets.AwardAssets             `json:"award_assets"`
}

//
// CREATE
//

type CreateDTO struct {
	Value                 int64 `json:"value" validate:"required,gt=0"`
	FileAchievementHeader *multipart.FileHeader
	FileAwardHeader       *multipart.FileHeader
	Description           *string `json:"description" validate:"omitempty,min=1"`
	Code                  string  `json:"code" validate:"required,min=1"`
	Name                  string  `json:"name" validate:"required,min=1"`
	ConditionType         string  `json:"condition_type" validate:"required,min=1"`
	Operator              string  `json:"operator" validate:"required,min=1"`
}

//
// CREATE ACHIEVEMENT
//

type CreateAchievementDTO struct {
	AchievementAssetsID int64   `json:"achievement_assets_id"`
	AwardAssetsID       int64   `json:"award_assets_id"`
	Description         *string `json:"description,omitempty"`
	Code                string  `json:"code"`
	Name                string  `json:"name"`
}

//
// CREATE ACHIEVEMENT CONDITION
//

type CreateAchievementConditionDTO struct {
	AchievementID int64  `json:"achievement_id"`
	Value         int64  `json:"value"`
	ConditionType string `json:"condition_type"`
	Operator      string `json:"operator"`
}

//
// SWAGGER
//

type DetailSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		Achievement       Achievement `json:"achievement"`
		Condition         Condition   `json:"condition"`
		AchievementAssets struct {
			ID                       int64     `json:"id"`
			Quality                  int       `json:"quality"`
			NameFile                 string    `json:"name_file"`
			NameFileWithoutExtension string    `json:"name_file_without_extension"`
			ServerPathFile           string    `json:"server_path_file"`
			ClientPathFile           string    `json:"client_path_file"`
			Extension                string    `json:"extension"`
			OldNameFile              string    `json:"old_name_file"`
			OldExtension             string    `json:"old_extension"`
			CreatedAt                time.Time `json:"created_at"`
			UpdatedAt                time.Time `json:"updated_at"`
		} `json:"achievement_assets"`
		AwardAssets struct {
			ID                       int64     `json:"id"`
			Quality                  int       `json:"quality"`
			NameFile                 string    `json:"name_file"`
			NameFileWithoutExtension string    `json:"name_file_without_extension"`
			ServerPathFile           string    `json:"server_path_file"`
			ClientPathFile           string    `json:"client_path_file"`
			Extension                string    `json:"extension"`
			OldNameFile              string    `json:"old_name_file"`
			OldExtension             string    `json:"old_extension"`
			CreatedAt                time.Time `json:"created_at"`
			UpdatedAt                time.Time `json:"updated_at"`
		} `json:"award_assets"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
