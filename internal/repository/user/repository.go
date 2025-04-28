package user

import "github.com/go-jedi/lingvogramm_backend/pkg/logger"

type Repository struct {
	logger logger.ILogger
}

func New(logger logger.ILogger) *Repository {
	return &Repository{
		logger: logger,
	}
}
