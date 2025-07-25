package middleware

import (
	"log"

	"github.com/go-jedi/lingramm_backend/config"
	adminguard "github.com/go-jedi/lingramm_backend/internal/middleware/admin_guard"
	"github.com/go-jedi/lingramm_backend/internal/middleware/auth"
	contentlengthlimiter "github.com/go-jedi/lingramm_backend/internal/middleware/content_length_limiter"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Middleware struct {
	Auth                 *auth.Middleware
	AdminGuard           *adminguard.Middleware
	ContentLengthLimiter *contentlengthlimiter.Middleware
}

func New(
	cfg config.MiddlewareConfig,
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
		Auth:                 auth.New(jwt, redis),
		AdminGuard:           adminguard.New(adminService, jwt),
		ContentLengthLimiter: contentlengthlimiter.New(cfg.ContentLengthLimiter.MaxBodySize),
	}
}
