package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	dailytask "github.com/go-jedi/lingramm_backend/internal/domain/daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/utils/nullify"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto dailytask.CreateDTO) (dailytask.DailyTask, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	r := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Create) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto dailytask.CreateDTO) (dailytask.DailyTask, error) {
	r.logger.Debug("[create a new daily task] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO daily_tasks(
		    words_learned_need,
		    tasks_completed_need,
		    lessons_finished_need,
		    words_translate_need,
		    dialog_completed_need,
		    experience_points_need,
		    is_active
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING *;
	`

	var ndt dailytask.DailyTask

	if err := tx.QueryRow(
		ctxTimeout, q,
		r.getArgs(dto)...,
	).Scan(
		&ndt.ID, &ndt.WordsLearnedNeed,
		&ndt.TasksCompletedNeed, &ndt.LessonsFinishedNeed,
		&ndt.WordsTranslateNeed, &ndt.DialogCompletedNeed,
		&ndt.ExperiencePointsNeed, &ndt.IsActive,
		&ndt.CreatedAt, &ndt.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new daily task", "err", err)
			return dailytask.DailyTask{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new daily task", "err", err)
		return dailytask.DailyTask{}, fmt.Errorf("could not create a new daily task: %w", err)
	}

	return ndt, nil
}

// getArgs get args.
func (r *Create) getArgs(dto dailytask.CreateDTO) []interface{} {
	return []interface{}{
		nullify.EmptyInt64WithDefault(dto.WordsLearnedNeed),
		nullify.EmptyInt64WithDefault(dto.TasksCompletedNeed),
		nullify.EmptyInt64WithDefault(dto.LessonsFinishedNeed),
		nullify.EmptyInt64WithDefault(dto.WordsTranslateNeed),
		nullify.EmptyInt64WithDefault(dto.DialogCompletedNeed),
		nullify.EmptyInt64WithDefault(dto.ExperiencePointsNeed),
		dto.IsActive,
	}
}
