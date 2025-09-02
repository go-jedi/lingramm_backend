package userstats

import "time"

type UserStats struct {
	ID                      int64      `json:"id"`
	StreakDays              int64      `json:"streak_days"`
	DailyTaskStreakDays     int64      `json:"daily_task_streak_days"`
	WordsLearned            int64      `json:"words_learned"`
	TasksCompleted          int64      `json:"tasks_completed"`
	LessonsFinished         int64      `json:"lessons_finished"`
	ExperienceScale         int64      `json:"experience_scale"`
	WordsTranslate          int64      `json:"words_translate"`
	DialoguesCompleted      int64      `json:"dialog_completed"`
	ExperiencePoints        int64      `json:"experience_points"`
	Level                   int64      `json:"level"`
	TelegramID              string     `json:"telegram_id"`
	LastStreakDay           time.Time  `json:"last_streak_day"`
	LastDailyTaskStreakDays *time.Time `json:"last_daily_task_streak_days,omitempty"`
	LastActiveAt            *time.Time `json:"last_active_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

//
// GET LEVEL BY TELEGRAM ID
//

type GetLevelByTelegramIDResponse struct {
	Level      int64  `json:"level"`
	TelegramID string `json:"telegram_id"`
}

//
// GET LEVEL INFO BY TELEGRAM ID
//

type GetLevelInfoByTelegramIDResponse struct {
	XPTotal       int64   `json:"xp_total"`
	Level         int64   `json:"level"`
	LevelFloorXP  int64   `json:"level_floor_xp"`
	LevelCeilXP   int64   `json:"level_ceil_xp"`
	NextLevel     int64   `json:"next_level"`
	XPInLevel     int64   `json:"xp_in_level"`
	XPLevelSize   int64   `json:"xp_level_size"`
	XPToNext      int64   `json:"xp_to_next"`
	ProgressRatio float64 `json:"progress_ratio"`
	LevelName     string  `json:"level_name"`
}

//
// SWAGGER
//

type GetLevelByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		Level      int64  `json:"level" example:"1"`
		TelegramID string `json:"telegram_id" example:"1"`
	} `json:"data"`
}

type GetLevelInfoByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		XPTotal       int64   `json:"xp_total" example:"0"`
		Level         int64   `json:"level" example:"1"`
		LevelFloorXP  int64   `json:"level_floor_xp" example:"0"`
		LevelCeilXP   int64   `json:"level_ceil_xp" example:"100"`
		NextLevel     int64   `json:"next_level" example:"2"`
		XPInLevel     int64   `json:"xp_in_level" example:"0"`
		XPLevelSize   int64   `json:"xp_level_size" example:"100"`
		XPToNext      int64   `json:"xp_to_next" example:"100"`
		ProgressRatio float64 `json:"progress_ratio" example:"0"`
		LevelName     string  `json:"level_name" example:"level 1"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
