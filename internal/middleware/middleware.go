package middleware

import (
	"log"

	adminguard "github.com/go-jedi/lingramm_backend/internal/middleware/admin_guard"
	"github.com/go-jedi/lingramm_backend/internal/middleware/auth"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Middleware struct {
	Auth       *auth.Middleware
	AdminGuard *adminguard.Middleware
}

func New(
	adminService *adminservice.Service,
	jwt *jwt.JWT,
	redis *redis.Redis,
) *Middleware {
	if jwt == nil {
		log.Fatal("jwt instance cannot be nil")
	}
	if redis == nil {
		log.Fatal("redis instance cannot be nil")
	}

	return &Middleware{
		Auth:       auth.New(jwt, redis),
		AdminGuard: adminguard.New(adminService, jwt),
	}
}
