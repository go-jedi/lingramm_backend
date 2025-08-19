package level

import "time"

// UserLevelHistory represents user level history in the system.
type UserLevelHistory struct {
	ID          int64     `json:"id"`
	LevelNumber int64     `json:"level_number"`
	XPAtReach   int64     `json:"xp_at_reach"`
	XPEventID   *int64    `json:"xp_event_id,omitempty"`
	TelegramID  string    `json:"telegram_id"`
	ReachedAt   time.Time `json:"reached_at"`
}

//
// CREATE USER LEVEL HISTORY DTO
//

type CreateUserLevelHistoryDTO struct {
	LevelNumber int64  `json:"level_number"`
	XPAtReach   int64  `json:"xp_at_reach"`
	XPEventID   *int64 `json:"xp_event_id,omitempty"`
	TelegramID  string `json:"telegram_id"`
}
