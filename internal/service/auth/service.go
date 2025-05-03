package auth

import (
	"github.com/go-jedi/lingvogramm_backend/internal/repository/user"
	signin "github.com/go-jedi/lingvogramm_backend/internal/service/auth/sign_in"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
)

type Service struct {
	SignIn signin.ISignIn
}

func New(
	userRepository *user.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	uuid uuid.IUUID,
) *Service {
	return &Service{
		SignIn: signin.New(userRepository, logger, postgres, uuid),
	}
}
