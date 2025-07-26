package auth

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/auth/check"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/auth/refresh"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/auth/sign_in"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	Check   check.ICheck
	Refresh refresh.IRefresh
	SignIn  signin.ISignIn
}

func New(
	userRepository *user.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	bigCache *bigcachepkg.BigCache,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		Check:   check.New(userRepository, logger, postgres, bigCache, jwt),
		Refresh: refresh.New(userRepository, logger, postgres, redis, bigCache, jwt),
		SignIn:  signin.New(userRepository, logger, postgres, redis, bigCache, jwt),
	}
}
