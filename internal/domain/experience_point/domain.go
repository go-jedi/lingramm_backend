package experiencepoint

import (
	"time"
)

type XPEvents struct {
	ID          int64     `json:"id"`
	EventTypeID int64     `json:"event_type_id"`
	TelegramID  string    `json:"telegram_id"`
	DeltaXP     int       `json:"delta_xp"`
	OccurredAt  time.Time `json:"occurred_at"`
	InsertedAt  time.Time `json:"inserted_at"`
	WeekStart   time.Time `json:"week_start"`
}

//
// CREATE XP EVENT
//

type CreateXPEventDTO struct {
	TelegramID string `json:"telegram_id"`
	EventType  string `json:"event_type"`
	DeltaXP    int64  `json:"delta_xp"`
}

//
// LEADERBOARD WEEKS PROCESS BATCH
//

type LeaderboardWeeksProcessBatchDTO struct {
	BatchSize          int64  `json:"batch_size"`
	StatementTimeoutMS int64  `json:"statement_timeout_ms"`
	LockTimeoutMS      int64  `json:"lock_timeout_ms"`
	WorkerName         string `json:"worker_name"`
}

type LeaderboardWeeksProcessBatchResponse struct {
	FromID         int64 `json:"from_id"`
	ToID           int64 `json:"to_id"`
	BatchCount     int64 `json:"batch_count"`
	NewEventCount  int64 `json:"new_event_count"`
	GroupsCount    int64 `json:"groups_count"`
	AppliedXP      int64 `json:"applied_xp"`
	NewLastEventID int64 `json:"new_last_event_id"`
	Processed      bool  `json:"processed"`
}

//
// GET LEADERBOARD TOP WEEK
//

type GetLeaderboardTopWeekDTO struct {
	Limit int64  `json:"limit" validate:"required,gt=0,lte=30"`
	TZ    string `json:"tz" validate:"required,oneof=Europe/Moscow"`
}

type GetLeaderboardTopWeekResponse struct {
	Position    int64  `json:"position"`
	XP          int64  `json:"xp"`
	TelegramID  string `json:"telegram_id"`
	DisplayName string `json:"display_name"`
	Medal       string `json:"medal"`
}

//
// GET LEADERBOARD TOP WEEK FOR USER
//

type GetLeaderboardTopWeekForUserDTO struct {
	Limit      int64  `json:"limit" validate:"required,gt=0,lte=30"`
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	TZ         string `json:"tz" validate:"required,oneof=Europe/Moscow"`
}

type GetLeaderboardTopWeekForUserResponse struct {
	Position    int64  `json:"position"`
	XP          int64  `json:"xp"`
	TelegramID  string `json:"telegram_id"`
	DisplayName string `json:"display_name"`
	Medal       string `json:"medal"`
}

//
// SWAGGER
//

type GetLeaderboardTopWeekSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		Position    int64  `json:"position" example:"1"`
		XP          int64  `json:"xp" example:"20"`
		TelegramID  string `json:"telegram_id" example:"1"`
		DisplayName string `json:"display_name" example:"some name"`
		Medal       string `json:"medal" example:"some medal"`
	} `json:"data"`
}

type GetLeaderboardTopWeekForUserSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		Position    int64  `json:"position" example:"1"`
		XP          int64  `json:"xp" example:"20"`
		TelegramID  string `json:"telegram_id" example:"1"`
		DisplayName string `json:"display_name" example:"some name"`
		Medal       string `json:"medal" example:"some medal"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
