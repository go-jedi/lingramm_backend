package experiencepoint

import "time"

type XPEvents struct {
	ID         int64     `json:"id"`
	TelegramID string    `json:"telegram_id"`
	DeltaXP    int       `json:"delta_xp"`
	OccurredAt time.Time `json:"occurred_at"`
	InsertedAt time.Time `json:"inserted_at"`
	Reason     string    `json:"reason"`
	WeekStart  time.Time `json:"week_start"`
}

//
// CREATE XP EVENT
//

type CreateXPEventDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Events     []struct {
		DeltaXP int    `json:"delta_xp" validate:"required,gte=-500,lte=500"`
		Reason  string `json:"reason" validate:"required,min=1,max=50"`
	} `json:"events" validate:"required,min=1,max=250"`
}

//
// LEADERBOARD WEEKS PROCESS BATCH
//

type LeaderboardWeeksProcessBatchDTO struct {
	WorkerName         string `json:"worker_name"`
	BatchSize          int64  `json:"batch_size"`
	StatementTimeoutMS int64  `json:"statement_timeout_ms"`
	LockTimeoutMS      int64  `json:"lock_timeout_ms"`
}

type LeaderboardWeeksProcessBatchResponse struct {
	Processed      bool  `json:"processed"`
	FromID         int64 `json:"from_id"`
	ToID           int64 `json:"to_id"`
	BatchCount     int64 `json:"batch_count"`
	NewEventCount  int64 `json:"new_event_count"`
	GroupsCount    int64 `json:"groups_count"`
	AppliedXP      int64 `json:"applied_xp"`
	NewLastEventID int64 `json:"new_last_event_id"`
}
