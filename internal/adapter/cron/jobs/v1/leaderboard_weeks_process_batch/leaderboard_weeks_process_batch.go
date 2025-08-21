package leaderboardweeksprocessbatch

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

const minTOAddMillisecond = 1500 // extra client-side headroom added to DB statement_timeout (in ms).

// LeaderboardWeeksProcessBatch periodically calls the DB function
// public.leaderboard_weeks_process_batch to fold xp_events into the weekly
// leaderboard aggregate. It runs in a burst: multiple back-to-back calls
// within a single tick until we catch up to the fixed "ceiling" (to_id),
// then sleeps until the next tick.
type LeaderboardWeeksProcessBatch struct {
	experiencePointService *experiencepointservice.Service
	logger                 *logger.Logger
	workerName             string
	batchSize              int64
	statementTimeoutMS     int64
	lockTimeoutMS          int64
	timeoutReliefCPU       int64
	sleepDuration          int
	timeout                int
}

// New constructs the cron job and starts it in a background goroutine.
// It does NOT block; call with a cancellable context to stop it later.
func New(
	ctx context.Context,
	experiencePointService *experiencepointservice.Service,
	cfg config.CronConfig,
	logger *logger.Logger,
) *LeaderboardWeeksProcessBatch {
	c := &LeaderboardWeeksProcessBatch{
		experiencePointService: experiencePointService,
		logger:                 logger,
		workerName:             cfg.LeaderboardWeeksProcessBatch.WorkerName,
		batchSize:              cfg.LeaderboardWeeksProcessBatch.BatchSize,
		statementTimeoutMS:     cfg.LeaderboardWeeksProcessBatch.StatementTimeoutMS,
		lockTimeoutMS:          cfg.LeaderboardWeeksProcessBatch.LockTimeoutMS,
		timeoutReliefCPU:       cfg.LeaderboardWeeksProcessBatch.TimeoutReliefCPU,
		sleepDuration:          cfg.LeaderboardWeeksProcessBatch.SleepDuration,
		timeout:                cfg.LeaderboardWeeksProcessBatch.Timeout,
	}

	go c.start(ctx)

	return c
}

// start sets up a ticker and on each tick performs a burst of DB function calls.
// The burst loops until the DB reports there is no more work within the fixed range,
// then waits for the next tick.
func (c *LeaderboardWeeksProcessBatch) start(ctx context.Context) {
	// tick every sleepDuration seconds (short cadence recommended in prod: 1–5s).
	ticker := time.NewTicker(time.Duration(c.sleepDuration) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("cron leaderboard weeks process batch stopped", slog.String("reason", ctx.Err().Error()))
			return
		case <-ticker.C:
			c.logger.Debug("[cron leaderboard weeks process batch] tick")

			// outer client timeout: must be slightly larger than DB statement_timeout.
			outerTO := time.Duration(c.timeout) * time.Second
			minTO := time.Duration(c.statementTimeoutMS+minTOAddMillisecond) * time.Millisecond
			if outerTO < minTO {
				outerTO = minTO
			}

			ctxTimeout, cancel := context.WithTimeout(ctx, outerTO)

			// run one burst (one or more function calls back-to-back).
			if err := c.processBatch(ctxTimeout); err != nil {
				// log but keep the cron alive; next tick will retry.
				c.logger.Error("error leaderboard weeks process batch", "err", err)
			}

			cancel()
		}
	}
}

// processBatch performs a burst:
// repeatedly calls the DB function once per iteration (one batch),
// until either: (a) no work is left right now, or (b) we've caught up
// to the fixed ceiling (new_last_event_id >= to_id).
func (c *LeaderboardWeeksProcessBatch) processBatch(ctx context.Context) error {
	var (
		data = experiencepoint.LeaderboardWeeksProcessBatchDTO{
			WorkerName:         c.workerName,
			BatchSize:          c.batchSize,
			StatementTimeoutMS: c.statementTimeoutMS,
			LockTimeoutMS:      c.lockTimeoutMS,
		}
		progress = false // whether any iteration actually processed data
		burst    = 0     // number of iterations in this burst (for observability)
	)

	for {
		// exactly one DB function call = one batch iteration.
		result, err := c.experiencePointService.LeaderboardWeeksProcessBatch.Execute(ctx, data)
		if err != nil {
			c.logger.Error("error execute leaderboard weeks process batch", "err", err)
			return err
		}

		// batch metrics for debugging/observability.
		c.logger.Debug("lbw batch",
			slog.Bool("processed", result.Processed),
			slog.Int64("from id", result.FromID),
			slog.Int64("to id", result.ToID),
			slog.Int64("batch count", result.BatchCount),
			slog.Int64("new event count", result.NewEventCount),
			slog.Int64("groups count", result.GroupsCount),
			slog.Int64("applied xp", result.AppliedXP),
			slog.Int64("new last event id", result.NewLastEventID),
		)

		if !result.Processed {
			// nothing to do at the moment — end the burst.
			break
		}

		progress = true
		burst++

		// if there's still a tail within the fixed ceiling (to_id),
		// immediately run the next iteration after a tiny CPU-friendly pause.
		if result.NewLastEventID < result.ToID {
			time.Sleep(time.Duration(c.timeoutReliefCPU) * time.Millisecond) // small pause to yield CPU.
			continue
		}

		// we've caught up to the ceiling for this burst — end the burst.
		break
	}

	if !progress {
		// not an error: simply no work right now.
		return nil
	}

	return nil
}
