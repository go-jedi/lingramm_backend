package undeletefileachievementcleaner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

var ErrLoadPathsFromRedis = errors.New("failed to load file paths from redis")

type UnDeleteFileAchievementCleaner struct {
	logger        *logger.Logger
	redis         *redis.Redis
	sleepDuration int
	timeout       int
}

func New(
	ctx context.Context,
	cfg config.CronConfig,
	logger *logger.Logger,
	redis *redis.Redis,
) *UnDeleteFileAchievementCleaner {
	c := &UnDeleteFileAchievementCleaner{
		logger:        logger,
		redis:         redis,
		sleepDuration: cfg.UnDeleteFileAchievementCleaner.SleepDuration,
		timeout:       cfg.UnDeleteFileAchievementCleaner.Timeout,
	}

	go c.Start(ctx)

	return c
}

func (c *UnDeleteFileAchievementCleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.sleepDuration) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("cron un delete file cleaner stopped", slog.String("reason", ctx.Err().Error()))
			return

		case <-ticker.C:
			c.logger.Debug("[un delete file cleaner] tick")

			ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.timeout)*time.Second)

			if err := c.cleanFiles(ctxTimeout); err != nil {
				c.logger.Error("error cleaning files", "err", err)
			}

			cancel()
		}
	}
}

func (c *UnDeleteFileAchievementCleaner) cleanFiles(ctx context.Context) error {
	paths, err := c.redis.UnDeleteFileAchievement.All(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLoadPathsFromRedis, err)
	}

	if len(paths) == 0 {
		c.logger.Debug("no files to delete")
		return nil
	}

	keysToDelete := make([]string, 0, len(paths))

	for k := range paths {
		p, ok := paths[k]
		if !ok {
			continue
		}

		if err := os.Remove(p); err != nil {
			c.logger.Warn("failed to remove file", "path", p, "error", err)
			continue
		}

		c.logger.Debug("file removed by cron", "path", p)

		keysToDelete = append(keysToDelete, k)
	}

	if len(keysToDelete) > 0 {
		if err := c.redis.UnDeleteFileAchievement.DeleteKeys(ctx, keysToDelete); err != nil {
			c.logger.Error("failed to delete keys from redis", "error", err)
		}
	}

	return nil
}
