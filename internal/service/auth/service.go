package auth

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/user"
	"github.com/go-jedi/lingramm_backend/internal/service/auth/check"
	signin "github.com/go-jedi/lingramm_backend/internal/service/auth/sign_in"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
)

type Service struct {
	SignIn signin.ISignIn
	Check  check.ICheck
}

func New(
	userRepository *user.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
	uuid uuid.IUUID,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		SignIn: signin.New(userRepository, logger, postgres, bigCache, uuid),
		Check:  check.New(userRepository, logger, postgres, bigCache, jwt),
	}
}
