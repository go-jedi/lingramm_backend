package middleware

import (
	"log"

	"github.com/go-jedi/lingramm_backend/config"
	adminguard "github.com/go-jedi/lingramm_backend/internal/middleware/admin_guard"
	"github.com/go-jedi/lingramm_backend/internal/middleware/auth"
	authwebsocket "github.com/go-jedi/lingramm_backend/internal/middleware/auth_websocket"
	contentlengthlimiter "github.com/go-jedi/lingramm_backend/internal/middleware/content_length_limiter"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Middleware struct {
	AdminGuard           *adminguard.Middleware
	Auth                 *auth.Middleware
	AuthWebSocket        *authwebsocket.Middleware
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
		AdminGuard:           adminguard.New(adminService, jwt),
		Auth:                 auth.New(jwt, redis),
		AuthWebSocket:        authwebsocket.New(jwt, redis),
		ContentLengthLimiter: contentlengthlimiter.New(cfg.ContentLengthLimiter.MaxBodySize),
	}
}
