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
