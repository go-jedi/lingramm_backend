package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	refreshtoken "github.com/go-jedi/lingramm_backend/pkg/redis/refresh_token"
	undeletefileachievement "github.com/go-jedi/lingramm_backend/pkg/redis/un_delete_file_achievement"
	undeletefileclient "github.com/go-jedi/lingramm_backend/pkg/redis/un_delete_file_client"
	"github.com/redis/go-redis/v9"
)

var ErrRedisPingFailed = errors.New("redis ping failed")

type Redis struct {
	RefreshToken            refreshtoken.IRefreshToken
	UnDeleteFileClient      undeletefileclient.IUnDeleteFileClient
	UnDeleteFileAchievement undeletefileachievement.IUnDeleteFileAchievement
}

func New(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	r := &Redis{}

	c := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		PoolTimeout:     time.Duration(cfg.PoolTimeout) * time.Second,
		TLSConfig:       nil,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: time.Duration(cfg.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(cfg.MaxRetryBackoff) * time.Millisecond,
	})

	if _, err := c.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisPingFailed, err)
	}

	r.RefreshToken = refreshtoken.New(cfg, c)
	r.UnDeleteFileClient = undeletefileclient.New(cfg.UnDeleteFileClient, c)
	r.UnDeleteFileAchievement = undeletefileachievement.New(cfg.UnDeleteFileAchievement, c)

	return r, nil
}
