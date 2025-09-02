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
	AchievementTypeID   int64     `json:"achievement_type_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Description         *string   `json:"description,omitempty"`
	Name                string    `json:"name"`
}

// Detail represents achievement detail in the system.
type Detail struct {
	Achievement       Achievement                         `json:"achievement"`
	AchievementAssets achievementassets.AchievementAssets `json:"achievement_assets"`
	AwardAssets       awardassets.AwardAssets             `json:"award_assets"`
}

//
// CREATE
//

type CreateDTO struct {
	FileAchievementHeader *multipart.FileHeader
	FileAwardHeader       *multipart.FileHeader
	Description           *string `json:"description" validate:"omitempty,min=1"`
	Name                  string  `json:"name" validate:"required,min=1"`
	AchievementType       string  `json:"achievement_type" validate:"required,min=1"`
}

//
// CREATE ACHIEVEMENT
//

type CreateAchievementDTO struct {
	AchievementAssetsID int64   `json:"achievement_assets_id"`
	AwardAssetsID       int64   `json:"award_assets_id"`
	AchievementTypeID   int64   `json:"achievement_type_id"`
	Description         *string `json:"description,omitempty"`
	Name                string  `json:"name"`
}

//
// SWAGGER
//

type AllDetailSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		Achievement struct {
			ID                  int64     `json:"id" example:"1"`
			AchievementAssetsID int64     `json:"achievement_assets_id"  example:"1"`
			AwardAssetsID       int64     `json:"award_assets_id"  example:"1"`
			AchievementTypeID   int64     `json:"achievement_type_id"  example:"3"`
			CreatedAt           time.Time `json:"created_at"  example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt           time.Time `json:"updated_at"  example:"2025-09-02T12:48:06.37622+03:00"`
			Description         *string   `json:"description,omitempty"  example:"description"`
			Name                string    `json:"name"  example:"two dialogs"`
		} `json:"achievement"`
		AchievementAssets struct {
			ID                       int64     `json:"id" example:"1"`
			Quality                  int       `json:"quality" example:"30"`
			NameFile                 string    `json:"name_file" example:"01K44X76FBXJYK4D153WHZFXH7.webp"`
			NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76FBXJYK4D153WHZFXH7"`
			ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/achievement/01K44X76FBXJYK4D153WHZFXH7.webp"`
			ClientPathFile           string    `json:"client_path_file" example:"/images/achievement/01K44X76FBXJYK4D153WHZFXH7.webp"`
			Extension                string    `json:"extension" example:".webp"`
			OldNameFile              string    `json:"old_name_file" example:"img.png"`
			OldExtension             string    `json:"old_extension" example:".png"`
			CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		} `json:"achievement_assets"`
		AwardAssets struct {
			ID                       int64     `json:"id" example:"1"`
			Quality                  int       `json:"quality" example:"30"`
			NameFile                 string    `json:"name_file" example:"01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76GAFBZBJ1W1WX4NSJT4"`
			ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/award/01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			ClientPathFile           string    `json:"client_path_file" example:"/images/award/01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			Extension                string    `json:"extension" example:".webp"`
			OldNameFile              string    `json:"old_name_file" example:"img.jpg"`
			OldExtension             string    `json:"old_extension" example:".jpg"`
			CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		} `json:"award_assets"`
	} `json:"data"`
}

type DetailSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		Achievement struct {
			ID                  int64     `json:"id" example:"1"`
			AchievementAssetsID int64     `json:"achievement_assets_id"  example:"1"`
			AwardAssetsID       int64     `json:"award_assets_id"  example:"1"`
			AchievementTypeID   int64     `json:"achievement_type_id"  example:"3"`
			CreatedAt           time.Time `json:"created_at"  example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt           time.Time `json:"updated_at"  example:"2025-09-02T12:48:06.37622+03:00"`
			Description         *string   `json:"description,omitempty"  example:"description"`
			Name                string    `json:"name"  example:"two dialogs"`
		} `json:"achievement"`
		AchievementAssets struct {
			ID                       int64     `json:"id" example:"1"`
			Quality                  int       `json:"quality" example:"30"`
			NameFile                 string    `json:"name_file" example:"01K44X76FBXJYK4D153WHZFXH7.webp"`
			NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76FBXJYK4D153WHZFXH7"`
			ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/achievement/01K44X76FBXJYK4D153WHZFXH7.webp"`
			ClientPathFile           string    `json:"client_path_file" example:"/images/achievement/01K44X76FBXJYK4D153WHZFXH7.webp"`
			Extension                string    `json:"extension" example:".webp"`
			OldNameFile              string    `json:"old_name_file" example:"img.png"`
			OldExtension             string    `json:"old_extension" example:".png"`
			CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		} `json:"achievement_assets"`
		AwardAssets struct {
			ID                       int64     `json:"id" example:"1"`
			Quality                  int       `json:"quality" example:"30"`
			NameFile                 string    `json:"name_file" example:"01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76GAFBZBJ1W1WX4NSJT4"`
			ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/award/01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			ClientPathFile           string    `json:"client_path_file" example:"/images/award/01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
			Extension                string    `json:"extension" example:".webp"`
			OldNameFile              string    `json:"old_name_file" example:"img.jpg"`
			OldExtension             string    `json:"old_extension" example:".jpg"`
			CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
			UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		} `json:"award_assets"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
