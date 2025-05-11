package middleware

import (
	adminguard "github.com/go-jedi/lingvogramm_backend/internal/middleware/admin_guard"
	"github.com/go-jedi/lingvogramm_backend/internal/middleware/auth"
	"github.com/go-jedi/lingvogramm_backend/pkg/jwt"
)

type Middleware struct {
	Auth       *auth.Middleware
	AdminGuard *adminguard.Middleware
}

func New(jwt *jwt.JWT) *Middleware {
	return &Middleware{
		Auth:       auth.New(jwt),
		AdminGuard: adminguard.New(jwt),
	}
}
