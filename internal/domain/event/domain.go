package event

//
// CREATE EVENT
//

type CreateEventsDTO struct {
	TelegramID string  `json:"telegram_id" validate:"required,min=1"`
	EventType  string  `json:"event_type" validate:"required,min=1,max=50"`
	Actions    Actions `json:"actions" validate:"required"`
}

type Actions struct {
	WordsLearned    *int64 `json:"words_learned,omitempty" validate:"omitempty,gt=0"`
	TasksCompleted  *int64 `json:"tasks_completed,omitempty" validate:"omitempty,gt=0"`
	LessonsFinished *int64 `json:"lessons_finished,omitempty" validate:"omitempty,gt=0"`
	WordsTranslate  *int64 `json:"words_translate,omitempty" validate:"omitempty,gt=0"`
	DialogCompleted *int64 `json:"dialog_completed,omitempty" validate:"omitempty,gt=0"`
}
