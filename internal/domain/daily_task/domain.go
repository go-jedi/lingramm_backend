package dailytask

import "time"

type DailyTask struct {
	ID                   int64     `json:"id"`
	WordsLearnedNeed     *int64    `json:"words_learned_need,omitempty"`
	TasksCompletedNeed   *int64    `json:"tasks_completed_need,omitempty"`
	LessonsFinishedNeed  *int64    `json:"lessons_finished_need,omitempty"`
	WordsTranslateNeed   *int64    `json:"words_translate_need,omitempty"`
	DialogCompletedNeed  *int64    `json:"dialog_completed_need,omitempty"`
	ExperiencePointsNeed *int64    `json:"experience_points_need,omitempty"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	WordsLearnedNeed     *int64 `json:"words_learned_need,omitempty" validate:"omitempty,gt=0"`
	TasksCompletedNeed   *int64 `json:"tasks_completed_need,omitempty" validate:"omitempty,gt=0"`
	LessonsFinishedNeed  *int64 `json:"lessons_finished_need,omitempty" validate:"omitempty,gt=0"`
	WordsTranslateNeed   *int64 `json:"words_translate_need,omitempty" validate:"omitempty,gt=0"`
	DialogCompletedNeed  *int64 `json:"dialog_completed_need,omitempty" validate:"omitempty,gt=0"`
	ExperiencePointsNeed *int64 `json:"experience_points_need,omitempty" validate:"omitempty,gt=0"`
	IsActive             bool   `json:"is_active"`
}
