package createtextcontent

import (
	"context"
	"errors"
	"fmt"
	"time"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateTextContent --output=mocks --case=underscore
type ICreateTextContent interface {
	Execute(ctx context.Context, tx pgx.Tx, dto localizedtext.CreateTextContentDTO) (localizedtext.TextContents, error)
}

type CreateTextContent struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateTextContent {
	r := &CreateTextContent{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateTextContent) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateTextContent) Execute(ctx context.Context, tx pgx.Tx, dto localizedtext.CreateTextContentDTO) (localizedtext.TextContents, error) {
	r.logger.Debug("[create text content] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO text_contents(
			code,
		    page, 
		    description
		) VALUES($1, $2, $3)
		RETURNING *;
	`

	var ntc localizedtext.TextContents

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.Code, dto.Page,
		dto.Description,
	).Scan(
		&ntc.ID, &ntc.Code,
		&ntc.Page, &ntc.Description,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create text content", "err", err)
			return localizedtext.TextContents{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create text content", "err", err)
		return localizedtext.TextContents{}, fmt.Errorf("could not create text content: %w", err)
	}

	return ntc, nil
}
