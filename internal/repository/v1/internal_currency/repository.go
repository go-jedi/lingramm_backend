package internalcurrency

import (
	adduserbalance "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency/add_user_balance"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency/get_user_balance"
	reduceuserbalance "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency/reduce_user_balance"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AddUserBalance    adduserbalance.IAddUserBalance
	GetUserBalance    getuserbalance.IGetUserBalance
	ReduceUserBalance reduceuserbalance.IReduceUserBalance
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AddUserBalance:    adduserbalance.New(queryTimeout, logger),
		GetUserBalance:    getuserbalance.New(queryTimeout, logger),
		ReduceUserBalance: reduceuserbalance.New(queryTimeout, logger),
	}
}
