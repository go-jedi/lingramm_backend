package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	ErrNatsURLIsEmpty                  = errors.New("nats url is empty")
	ErrStreamNameIsEmpty               = errors.New("stream name is empty")
	ErrSubjectIsEmpty                  = errors.New("subject is empty")
	ErrTimeoutIsEmpty                  = errors.New("timeout is empty")
	ErrNameInStreamOptionIsEmpty       = errors.New("name in stream option is empty")
	ErrNameInNatsOptionIsEmpty         = errors.New("name in nats option is empty")
	ErrDurableInSubscribeOptionIsEmpty = errors.New("durable in subscribe option is empty")
	ErrNatsConnect                     = errors.New("nats connect error")
	ErrGettingJetStream                = errors.New("error getting JetStream")
)

// Consumer is responsible for subscribing to a JetStream subject,
// handling messages, and delivering them to an external handler function.
type Consumer struct {
	nc         *nats.Conn            // Low-level NATS connection.
	js         nats.JetStreamContext // JetStream context for stream and consumer APIs.
	sub        *nats.Subscription    // Durable subscription to receive messages.
	opts       options               // Configuration values for NATS, stream, and subscription.
	closeOnce  sync.Once             // Ensures shutdown logic is executed only once.
	telegramID string                // Used to generate unique subject/durable names per user.
}

func New(cfg config.NatsNotificationConfig, telegramID string) (*Consumer, error) {
	c := &Consumer{
		opts:       getOptions(cfg),
		telegramID: telegramID,
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	nc, err := nats.Connect(c.opts.natsURL, c.getNatsOptions()...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNatsConnect, err)
	}
	c.nc = nc

	js, err := nc.JetStream()
	if err != nil {
		c.nc.Close()
		return nil, fmt.Errorf("%w: %v", ErrGettingJetStream, err)
	}
	c.js = js

	return c, nil
}

// validate ensures all required config parameters are present before starting the consumer.
// This avoids runtime panics or misbehavior due to missing fields.
func (c *Consumer) validate() error {
	if c.opts.natsURL == "" {
		return ErrNatsURLIsEmpty
	}
	if c.opts.streamName == "" {
		return ErrStreamNameIsEmpty
	}
	if c.opts.subject == "" {
		return ErrSubjectIsEmpty
	}
	if c.opts.timeout == 0 {
		return ErrTimeoutIsEmpty
	}
	if c.opts.streamOption.name == "" {
		return ErrNameInStreamOptionIsEmpty
	}
	if c.opts.natsOption.name == "" {
		return ErrNameInNatsOptionIsEmpty
	}
	if c.opts.subscribeOption.durable == "" {
		return ErrDurableInSubscribeOptionIsEmpty
	}

	return nil
}

// Start creates a durable subscription to a JetStream subject.
// Messages are handled asynchronously using a provided handler function.
func (c *Consumer) Start(ctx context.Context, fn func(notification.Notification) error) error {
	subject := c.opts.subject + c.telegramID
	durable := c.opts.subscribeOption.durable + c.telegramID

	if err := c.ensureStream(); err != nil {
		return err
	}

	log.Printf("subscribing to [%s] with durable [%s]", subject, durable)

	sub, err := c.js.Subscribe(subject, c.messageHandler(ctx, fn), c.getSubscribeOptions()...)
	if err != nil {
		log.Printf("subscription error [%s]: %v", subject, err)
		return err
	}

	c.sub = sub

	return nil
}

// messageHandler returns a function that:
// 1. Decodes the message (msgpack),
// 2. Calls the handler within a timeout context,
// 3. Acks the message on success,
// 4. NAKs on failure or timeout,
// 5. Terminates the message on decode error or panic.
// This ensures safe and reliable message processing.
func (c *Consumer) messageHandler(ctx context.Context, fn func(notification.Notification) error) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in message handler: %v", r)
				if err := msg.Term(); err != nil {
					log.Printf("error Term(): %v", err)
				}
			}
		}()

		var n notification.Notification
		if err := msgpack.Unmarshal(msg.Data, &n); err != nil {
			log.Printf("failed to decode msgpack: %v", err)
			if err := msg.Term(); err != nil {
				log.Printf("error Term(): %v", err)
			}
			return
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.opts.timeout)*time.Second)
		defer cancel()

		errCh := make(chan error, 1)
		go func() {
			defer close(errCh)
			errCh <- fn(n)
		}()

		select {
		case <-ctxTimeout.Done():
			log.Printf("handler timeout for subject [%s]", msg.Subject)
			if err := msg.Nak(); err != nil {
				log.Printf("error Nak(): %v", err)
			}
		case err := <-errCh:
			if err != nil {
				log.Printf("external handler error: %v", err)
				if err := msg.Nak(); err != nil {
					log.Printf("error Nak(): %v", err)
				}
			} else {
				if err := msg.Ack(); err != nil {
					log.Printf("⚠️ Ack error: %v", err)
				}
			}
		}
	}
}

// Close gracefully shuts down the NATS connection and subscription.
// Ensures in-flight messages are drained before exit.
func (c *Consumer) Close(ctx context.Context) error {
	log.Println("terminating NATS connection...")

	var err error
	c.closeOnce.Do(func() {
		if c.sub != nil {
			if e := c.sub.Drain(); e != nil {
				log.Printf("subscription drain error: %v", e)
				err = e
			}
		}
		if c.nc != nil && !c.nc.IsClosed() {
			if e := c.nc.Drain(); e != nil {
				log.Printf("connection drain error: %v", e)
				err = e
			}
		}
	})

	const timeout = 5

	select {
	case <-ctx.Done():
	case <-time.After(time.Duration(timeout) * time.Second):
		log.Println("force exit after timeout")
	}

	log.Println("connection closed")

	return err
}

// ensureStream checks if the stream exists, and creates it if missing.
// This allows the consumer to auto-recover on cold starts or fresh environments.
func (c *Consumer) ensureStream() error {
	_, err := c.js.StreamInfo(c.opts.streamName)
	if err == nil {
		return nil
	}

	log.Printf("creating new stream [%s]...", c.opts.streamName)

	_, err = c.js.AddStream(c.getStreamOptions())
	if err != nil {
		log.Printf("failed to add stream: %v", err)
		return err
	}

	return nil
}

// getNatsOptions builds all connection-related options for NATS, including
// reconnect behavior, jitter, ping settings, and logging hooks.
func (c *Consumer) getNatsOptions() []nats.Option {
	opts := []nats.Option{
		nats.MaxReconnects(c.opts.natsOption.maxReconnects),
		nats.ReconnectWait(c.opts.natsOption.reconnectWait),
		nats.ReconnectJitter(
			c.opts.natsOption.reconnectJitter.jitter,
			c.opts.natsOption.reconnectJitter.jitterForTLS,
		),
		nats.Timeout(c.opts.natsOption.timeout),
		nats.DrainTimeout(c.opts.natsOption.drainTimeout),
		nats.PingInterval(c.opts.natsOption.pingInterval),
		nats.MaxPingsOutstanding(c.opts.natsOption.maxPingsOutstanding),
		nats.Name(c.opts.natsOption.name + c.telegramID),
	}

	if c.opts.natsOption.closedHandler {
		opts = append(opts, nats.ClosedHandler(func(_ *nats.Conn) {
			log.Println("NATS connection closed")
		}))
	}
	if c.opts.natsOption.reconnectHandler {
		opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("reconnected to: %s", nc.ConnectedUrl())
		}))
	}
	if c.opts.natsOption.disconnectErrHandler {
		opts = append(opts, nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("lost connection: %v", err)
		}))
	}
	if c.opts.natsOption.errorHandler {
		opts = append(opts, nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("async error: %v", err)
		}))
	}

	return opts
}

// getStreamOptions returns stream settings that must match the publisher.
// If stream options differ between publisher and consumer, JetStream behavior may become inconsistent.
func (c *Consumer) getStreamOptions() *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:      c.opts.streamName,
		Subjects:  []string{c.opts.subject + c.telegramID},
		Storage:   nats.StorageType(c.opts.streamOption.storage),
		Retention: nats.RetentionPolicy(c.opts.streamOption.retention),
		MaxAge:    c.opts.streamOption.maxAge,
		MaxMsgs:   c.opts.streamOption.maxMsgs,
		MaxBytes:  c.opts.streamOption.maxBytes,
		Discard:   nats.DiscardPolicy(c.opts.streamOption.discard),
	}
}

// getSubscribeOptions returns subscription options:
// - Durable: stores delivery state for resuming after reconnect,
// - AckWait: controls redelivery time,
// - MaxAckPending: limits number of in-flight messages.
func (c *Consumer) getSubscribeOptions() []nats.SubOpt {
	opts := []nats.SubOpt{
		nats.Durable(c.opts.subscribeOption.durable + c.telegramID),
		nats.AckWait(c.opts.subscribeOption.ackWait),
		nats.MaxAckPending(c.opts.subscribeOption.maxAckPending),
	}

	if c.opts.subscribeOption.manualAck {
		opts = append(opts, nats.ManualAck())
	}
	if c.opts.subscribeOption.deliverAll {
		opts = append(opts, nats.DeliverAll())
	}

	return opts
}
