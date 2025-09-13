package stream

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	notificationhub "github.com/go-jedi/lingramm_backend/pkg/ws_manager/notification"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

const (
	readDeadline          = 90 * time.Second // how long to wait before timing out a read.
	writeTimeout          = 5 * time.Second  // max duration to write to a client.
	pingInterval          = 30 * time.Second // how often to send a ping frame to the client.
	wsCloseWriteTimeout   = 2 * time.Second  // timeout for sending close control frame.
	liveChannelBufferSize = 256              // buffered channel size for live notifications.
	wsCloseReason         = "bye"            // close frame reason.
)

// Stream manages WebSocket connections and notification streaming.
type Stream struct {
	notificationService  *notificationservice.Service
	userStatsService     *userstatsservice.Service
	userDailyTaskService *userdailytaskservice.Service
	logger               logger.ILogger
	rabbitMQ             *rabbitmq.RabbitMQ
	redis                *redis.Redis
	middleware           *middleware.Middleware
	hub                  *notificationhub.Hub
}

// New returns a new instance of Stream.
func New(
	notificationService *notificationservice.Service,
	userStatsService *userstatsservice.Service,
	userDailyTaskService *userdailytaskservice.Service,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	redis *redis.Redis,
	middleware *middleware.Middleware,
	hub *notificationhub.Hub,
) *Stream {
	return &Stream{
		notificationService:  notificationService,
		userStatsService:     userStatsService,
		userDailyTaskService: userDailyTaskService,
		logger:               logger,
		rabbitMQ:             rabbitMQ,
		redis:                redis,
		middleware:           middleware,
		hub:                  hub,
	}
}

// Execute opens a WebSocket stream for user notifications.
// @Summary Notifications WebSocket stream
// @Description Upgrades the connection to WebSocket and streams notifications for the specified Telegram ID.
// @Description Server sends periodic pings; client may send `{"type":"ACK","id":<notification_id>}` to confirm delivery
// @Description and `{"type":"PONG"}` to refresh presence.
// @Tags Notification
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Success 101 {string} string "Switching Protocols (WebSocket established)"
// @Failure 400 {object} notification.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} notification.ErrorSwaggerResponse "Internal server error"
// @Router /v1/ws/notification/stream [get]
func (h *Stream) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get notifications stream] execute handler")

	fmt.Println("WS handler reached:", c.Path())
	fmt.Println("4")

	telegramID, err := h.middleware.AuthWebSocket.GetTelegramIDFromContext(c)
	if err != nil {
		h.logger.Error("failed to get telegramID", "error", "failed to get telegramID")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get telegramID", "failed to get telegramID", nil))
	}

	fmt.Println("5")

	return h.upgradeAndServe(c, telegramID)
}

func (h *Stream) upgradeAndServe(c fiber.Ctx, telegramID string) error {
	u := websocket.FastHTTPUpgrader{
		CheckOrigin: func(_ *fasthttp.RequestCtx) bool { return true },
	}

	return u.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
		h.runSession(conn, telegramID)
	})
}

func (h *Stream) runSession(conn *websocket.Conn, telegramID string) {
	defer h.closeWS(conn)

	// cancelable background context for the WS session lifetime.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// register the connection in the Hub for this telegram id.
	ce := h.registerInHub(telegramID, conn, cancel)
	defer h.hub.Delete(telegramID)

	// initialize online presence & pong behavior.
	if err := h.initPresence(ctx, conn, telegramID); err != nil {
		h.logger.Error("init user presence failed", "error", err)
		return
	}

	// send ping frames periodically to check if the connection is alive.
	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	// t0 marks the point in time separating old and new notifications.
	t0 := time.Now().UTC()

	// channels for live notifications, errors, and completion.
	liveMsgCh := make(chan notification.SendNotificationDTO, liveChannelBufferSize)
	liveErrCh := make(chan error, 1)
	liveDone := make(chan struct{})
	readErrCh := make(chan error, 1)

	// start background workers.
	go h.consumeNotifications(ctx, telegramID, liveMsgCh, liveErrCh, liveDone) // RabbitMQ live feed.
	go h.getAllPendingNotificationsFromDB(ctx, ce, telegramID, t0)             // pending notifications.
	go h.ensureStreakDaysIncrementToday(ctx, telegramID)                       // ensure streak days increment today.
	go h.assignDailyTask(ctx, telegramID)                                      // assign daily task.
	go h.startReadLoop(ctx, conn, telegramID, readErrCh)                       // client messages (ACK/PONG).

	// main loop to handle pinging, sending live notifications, and errors.
	h.runEventLoop(ctx, ce, pingTicker, liveMsgCh, liveErrCh, readErrCh, liveDone)
}

// registerInHub register connection in hub.
func (h *Stream) registerInHub(telegramID string, conn *websocket.Conn, cancel context.CancelFunc) *notificationhub.ConnectionEntry {
	ce := &notificationhub.ConnectionEntry{
		Connection: conn,
		Cancel:     cancel,
	}
	h.hub.Set(telegramID, ce)
	return ce
}

func (h *Stream) initPresence(ctx context.Context, conn *websocket.Conn, telegramID string) error {
	// set initial online status with TTL.
	if err := h.redis.UserPresence.Set(ctx, telegramID); err != nil {
		h.logger.Error("redis set online failed", "error", err)
		return err
	}

	// pong handler: refresh TTL on pong frame from client.
	conn.SetPongHandler(func(_ string) error {
		ok, err := h.redis.UserPresence.RefreshTTL(ctx, telegramID)
		if err != nil {
			h.logger.Error("redis refresh user presence ttl failed (pong handler)", "error", err)
			// Do not close the socket due to a transient Redis error.
			return nil
		}
		if !ok {
			if err := h.redis.UserPresence.Set(ctx, telegramID); err != nil {
				h.logger.Error("redis set user presence after missing key failed (pong handler)", "error", err)
			}
		}
		return nil
	})

	return nil
}

// runEventLoop loop to handle pinging, sending live notifications, and errors.
func (h *Stream) runEventLoop(
	ctx context.Context,
	ce *notificationhub.ConnectionEntry,
	pingTicker *time.Ticker,
	liveMsgCh <-chan notification.SendNotificationDTO,
	liveErrCh <-chan error,
	readErrCh <-chan error,
	liveDone <-chan struct{},
) {
	for {
		select {
		case <-pingTicker.C:
			// send ping frame under the same mutex as JSON writes to avoid concurrent writes
			if err := h.writeControl(ce, websocket.PingMessage, nil); err != nil {
				h.logger.Warn("failed to send ping", "error", err)
				return
			}
		case n := <-liveMsgCh:
			// send live notification to client.
			if err := h.sendJSON(ce, n); err != nil {
				h.logger.Error("send live notification failed", "err", err)
				return
			}
		case err := <-liveErrCh:
			// error from live consumer.
			if err != nil && !errors.Is(err, context.Canceled) {
				h.logger.Error("live consumer error", "err", err)
				return
			}
		case err := <-readErrCh:
			// error from client message loop (connection closed or failure).
			if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
				h.logger.Warn("read loop ended", "err", err)
			}

			return
		case <-liveDone:
			// live consumer finished.
			return
		case <-ctx.Done():
			// context canceled (connection closed).
			return
		}
	}
}

// closeWS tries to close WebSocket with a close frame before disconnecting.
func (h *Stream) closeWS(conn *websocket.Conn) {
	if conn == nil {
		return
	}

	if err := conn.SetWriteDeadline(time.Now().Add(wsCloseWriteTimeout)); err != nil {
		h.logger.Warn("set write deadline before close failed", "error", err)
	}

	if err := conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, wsCloseReason),
		time.Now().Add(wsCloseWriteTimeout),
	); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		h.logger.Warn("write close control failed", "error", err)
	}

	if err := conn.Close(); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		h.logger.Warn("websocket close failed", "error", err)
	}
}

// sendJSON sends a JSON message to the client with write deadline.
func (h *Stream) sendJSON(ce *notificationhub.ConnectionEntry, v interface{}) error {
	ce.Mu.Lock()
	defer ce.Mu.Unlock()

	if err := ce.Connection.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}

	return ce.Connection.WriteJSON(v)
}

// writeControl writes a control frame (e.g., ping) under the same mutex to avoid concurrent writes.
func (h *Stream) writeControl(ce *notificationhub.ConnectionEntry, mt int, data []byte) error {
	ce.Mu.Lock()
	defer ce.Mu.Unlock()

	if err := ce.Connection.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}

	return ce.Connection.WriteControl(mt, data, time.Now().Add(writeTimeout))
}

// consumeNotifications listens to RabbitMQ for new notifications in real-time.
// It sends them to liveMsgCh and reports errors to liveErrCh.
func (h *Stream) consumeNotifications(
	ctx context.Context,
	telegramID string,
	liveMsgCh chan<- notification.SendNotificationDTO,
	liveErrCh chan<- error,
	liveDone chan<- struct{},
) {
	defer close(liveDone)

	err := h.rabbitMQ.Notification.Consumer.Execute(ctx, telegramID, func(msg notification.SendNotificationDTO) error {
		select {
		case liveMsgCh <- msg:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		select {
		case liveErrCh <- err:
		default:
		}
	}
}

// ensureStreakDaysIncrementToday ensure streak days increment today.
func (h *Stream) ensureStreakDaysIncrementToday(ctx context.Context, telegramID string) {
	// ensure streak days increment.
	if err := h.userStatsService.EnsureStreakDaysIncrementToday.Execute(ctx, telegramID); err != nil {
		h.logger.Warn("ensure streak days increment today failed", "error", err)
	}
}

// assignDailyTask assign daily task.
func (h *Stream) assignDailyTask(ctx context.Context, telegramID string) {
	// check assign daily task exists by telegram id.
	isAssignDailyTask, err := h.userDailyTaskService.ExistsAssignDailyTaskByTelegramID.Execute(ctx, telegramID)
	if err != nil {
		h.logger.Warn("check assign daily task exists failed", "error", err)
	}

	if !isAssignDailyTask { // if daily task does not assign.
		// assign daily task by telegram id.
		if _, err := h.userDailyTaskService.AssignDailyTaskByTelegramID.Execute(ctx, telegramID); err != nil {
			h.logger.Warn("assign daily task failed", "error", err)
		}
	}
}

// getAllPendingNotificationsFromDB gets all pending notifications (<= t0) by Telegram ID and sends them to the client.
func (h *Stream) getAllPendingNotificationsFromDB(
	ctx context.Context,
	ce *notificationhub.ConnectionEntry,
	telegramID string,
	t0 time.Time,
) {
	result, err := h.notificationService.AllPendingBeforeByTelegramID.Execute(ctx, telegramID, t0)
	if err != nil {
		h.logger.Warn("get all pending before by telegramID failed", "error", err)
		return
	}

	for i := range result {
		n := notification.SendNotificationDTO{
			ID:         result[i].ID,
			Message:    result[i].Message,
			Type:       result[i].Type,
			TelegramID: result[i].TelegramID,
			CreatedAt:  result[i].CreatedAt,
		}
		if err := h.sendJSON(ce, n); err != nil {
			h.logger.Error("send pending notification failed", "error", err)
			return
		}
	}
}

// startReadLoop continuously reads messages from the client.
// It handles ACKs and PONGs (refreshes TTL on PONG with fallback).
func (h *Stream) startReadLoop(
	ctx context.Context,
	conn *websocket.Conn,
	telegramID string,
	readErrCh chan<- error,
) {
	defer close(readErrCh)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(readDeadline)); err != nil {
			readErrCh <- err
			return
		}

		mt, p, err := conn.ReadMessage()
		if err != nil {
			readErrCh <- err
			return
		}

		if mt != websocket.TextMessage && mt != websocket.BinaryMessage {
			continue
		}

		var msg notification.ACKMessage
		if err := jsoniter.Unmarshal(p, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "ACK":
			if msg.ID == 0 {
				continue
			}

			if err := h.notificationService.UpdateStatus.Execute(ctx, msg.ID, notification.SentStatus); err != nil {
				h.logger.Error("mark sent failed", "id", msg.ID, "error", err)
			}
		case "PONG":
			ok, err := h.redis.UserPresence.RefreshTTL(ctx, telegramID)
			if err != nil {
				h.logger.Error("redis refresh user presence ttl failed", "error", err)
				break
			}
			if !ok {
				if err := h.redis.UserPresence.Set(ctx, telegramID); err != nil {
					h.logger.Error("redis set user presence after missing key failed", "error", err)
				}
			}
		}
	}
}
