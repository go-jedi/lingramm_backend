package userdailytask

import (
	"time"
)

type Requirements struct {
	WordsLearnedNeed     *int64 `json:"words_learned_need,omitempty"`
	TasksCompletedNeed   *int64 `json:"tasks_completed_need,omitempty"`
	LessonsFinishedNeed  *int64 `json:"lessons_finished_need,omitempty"`
	WordsTranslateNeed   *int64 `json:"words_translate_need,omitempty"`
	DialogCompletedNeed  *int64 `json:"dialog_completed_need,omitempty"`
	ExperiencePointsNeed *int64 `json:"experience_points_need,omitempty"`
}

type Progress struct {
	WordsLearned     *int64 `json:"words_learned,omitempty"`
	TasksCompleted   *int64 `json:"tasks_completed,omitempty"`
	LessonsFinished  *int64 `json:"lessons_finished,omitempty"`
	WordsTranslate   *int64 `json:"words_translate,omitempty"`
	DialogCompleted  *int64 `json:"dialog_completed,omitempty"`
	ExperiencePoints *int64 `json:"experience_points,omitempty"`
}

type ProgressPercent struct {
	WordsLearned     *int64 `json:"words_learned,omitempty"`
	TasksCompleted   *int64 `json:"tasks_completed,omitempty"`
	LessonsFinished  *int64 `json:"lessons_finished,omitempty"`
	WordsTranslate   *int64 `json:"words_translate,omitempty"`
	DialogCompleted  *int64 `json:"dialog_completed,omitempty"`
	ExperiencePoints *int64 `json:"experience_points,omitempty"`
}

type Actions struct {
	WordsLearned     *int64 `json:"words_learned,omitempty"`
	TasksCompleted   *int64 `json:"tasks_completed,omitempty"`
	LessonsFinished  *int64 `json:"lessons_finished,omitempty"`
	WordsTranslate   *int64 `json:"words_translate,omitempty"`
	DialogCompleted  *int64 `json:"dialog_completed,omitempty"`
	ExperiencePoints *int64 `json:"experience_points,omitempty"`
}

//
// ASSIGN DAILY TASK BY TELEGRAM ID
//

type AssignDailyTaskByTelegramIDResponse struct {
	ID              int64           `json:"id"`
	Date            time.Time       `json:"date"`
	IsCompleted     bool            `json:"is_completed"`
	Requirements    Requirements    `json:"requirements"`
	Progress        Progress        `json:"progress"`
	ProgressPercent ProgressPercent `json:"progress_percent"`
}

// ConvertToGetCurrentDailyTask convert to get current daily task.
func (adt *AssignDailyTaskByTelegramIDResponse) ConvertToGetCurrentDailyTask() GetCurrentDailyTaskByTelegramIDResponse {
	return GetCurrentDailyTaskByTelegramIDResponse{
		ID:              adt.ID,
		Date:            adt.Date,
		IsCompleted:     adt.IsCompleted,
		Requirements:    adt.Requirements,
		Progress:        adt.Progress,
		ProgressPercent: adt.ProgressPercent,
	}
}

//
// GET CURRENT DAILY TASK BY TELEGRAM ID
//

type GetCurrentDailyTaskByTelegramIDResponse struct {
	ID              int64           `json:"id"`
	Date            time.Time       `json:"date"`
	IsCompleted     bool            `json:"is_completed"`
	Requirements    Requirements    `json:"requirements"`
	Progress        Progress        `json:"progress"`
	ProgressPercent ProgressPercent `json:"progress_percent"`
}

//
// SYNC USER DAILY TASK PROGRESS
//

type SyncUserDailyTaskProgressDTO struct {
	TelegramID string  `json:"telegram_id"`
	Actions    Actions `json:"actions"`
}
