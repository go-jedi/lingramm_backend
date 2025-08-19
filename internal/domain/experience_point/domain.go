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
