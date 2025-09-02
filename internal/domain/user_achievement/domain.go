package userachievement

import "time"

// Detail represents user achievement detail in the system.
type Detail struct {
	ID                  int64  `json:"id"`
	TelegramID          string `json:"telegram_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	AchievementPathFile string `json:"achievement_path_file"`
	AwardPathFile       string `json:"award_path_file"`
}

type UnlockAvailableAchievementsResponse struct {
	AchievementID   int64     `json:"achievement_id"`
	AchievementName string    `json:"achievement_name"`
	UnlockedAt      time.Time `json:"unlocked_at"`
}

//
// SWAGGER
//

type AllDetailByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID                  int64  `json:"id" example:"1"`
		TelegramID          string `json:"telegram_id" example:"1"`
		Name                string `json:"name" example:"some name"`
		Description         string `json:"description" example:"some description"`
		AchievementPathFile string `json:"achievement_path_file" example:"/images/achievement/01K44X76FBXJYK4D153WHZFXH7.webp"`
		AwardPathFile       string `json:"award_path_file" example:"/images/award/01K44X76GAFBZBJ1W1WX4NSJT4.webp"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
