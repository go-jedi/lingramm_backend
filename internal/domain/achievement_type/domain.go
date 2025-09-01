package achievementtype

import "time"

type AchievementType struct {
	ID                      int64     `json:"id"`
	Name                    string    `json:"name"`
	Description             *string   `json:"description,omitempty"`
	StreakDaysNeed          *int64    `json:"streak_days_need,omitempty"`
	DailyTaskStreakDaysNeed *int64    `json:"daily_task_streak_days_need,omitempty"`
	WordsLearnedNeed        *int64    `json:"words_learned_need,omitempty"`
	TasksCompletedNeed      *int64    `json:"tasks_completed_need,omitempty"`
	LessonsFinishedNeed     *int64    `json:"lessons_finished_need,omitempty"`
	WordsTranslateNeed      *int64    `json:"words_translate_need,omitempty"`
	DialogCompletedNeed     *int64    `json:"dialog_completed_need,omitempty"`
	ExperiencePointsNeed    *int64    `json:"experience_points_need,omitempty"`
	LevelNeed               *int64    `json:"level_need,omitempty"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
