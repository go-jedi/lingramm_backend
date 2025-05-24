package internalcurrency

import (
	internalcurrency "github.com/go-jedi/lingramm_backend/internal/repository/internal_currency"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"
	getuserbalance "github.com/go-jedi/lingramm_backend/internal/service/internal_currency/get_user_balance"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	GetUserBalance getuserbalance.IGetUserBalance
}

func New(
	internalCurrencyRepository *internalcurrency.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *Service {
	return &Service{
		GetUserBalance: getuserbalance.New(
			internalCurrencyRepository,
			userRepository,
			logger,
			postgres,
			bigCache,
		),
	}
}
