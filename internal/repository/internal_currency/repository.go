package internalcurrency

import (
	getuserbalance "github.com/go-jedi/lingramm_backend/internal/repository/internal_currency/get_user_balance"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	GetUserBalance getuserbalance.IGetUserBalance
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		GetUserBalance: getuserbalance.New(queryTimeout, logger),
	}
}
