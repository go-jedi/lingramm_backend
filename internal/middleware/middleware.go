package middleware

import (
	"log"

	adminguard "github.com/go-jedi/lingvogramm_backend/internal/middleware/admin_guard"
	"github.com/go-jedi/lingvogramm_backend/internal/middleware/auth"
	adminservice "github.com/go-jedi/lingvogramm_backend/internal/service/admin"
	"github.com/go-jedi/lingvogramm_backend/pkg/jwt"
)

type Middleware struct {
	Auth       *auth.Middleware
	AdminGuard *adminguard.Middleware
}

func New(
	adminService *adminservice.Service,
	jwt *jwt.JWT,
) *Middleware {
	if jwt == nil {
		log.Fatal("JWT instance cannot be nil")
	}

	return &Middleware{
		Auth:       auth.New(jwt),
		AdminGuard: adminguard.New(adminService, jwt),
	}
}
