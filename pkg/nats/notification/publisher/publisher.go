package publisher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

// Error definitions for better readability and error handling.
var (
	ErrNatsURLIsEmpty                   = errors.New("nats url is empty")
	ErrStreamNameIsEmpty                = errors.New("stream name is empty")
	ErrSubjectIsEmpty                   = errors.New("subject is empty")
	ErrTimeoutIsEmpty                   = errors.New("timeout is empty")
	ErrNameInStreamOptionIsEmpty        = errors.New("name in stream option is empty")
	ErrNameInNatsOptionIsEmpty          = errors.New("name in nats option is empty")
	ErrNatsConnect                      = errors.New("nats connect error")
	ErrGettingJetStream                 = errors.New("error getting JetStream")
	ErrEnsureStreamFailed               = errors.New("ensureStream failed")
	ErrFailedToMarshalNotification      = errors.New("failed to marshal notification")
	ErrPublishAsyncTimeoutForSubject    = errors.New("publish async timeout for subject")
	ErrAckFutureOkReturnedNil           = errors.New("ackFuture.Ok returned nil")
	ErrReceivedAckForUnexpectedStream   = errors.New("received ack for unexpected stream")
	ErrPublishAckFailed                 = errors.New("publish ack failed")
	ErrTimeoutWaitingForAsyncCompletion = errors.New("timeout waiting for async completion")
)

// IPublisher defines the interface for the publisher nats.
//
//go:generate mockery --name=IPublisher --output=mocks --case=underscore
type IPublisher interface {
	Start(ctx context.Context, telegramID string, data notification.Notification) error
	Close() error
	WaitForCompletion(ctx context.Context) error
}

// Publisher handles connection to NATS and JetStream publishing logic.
type Publisher struct {
	nc   *nats.Conn            // NATS connection.
	js   nats.JetStreamContext // JetStream context.
	opts options               // configuration options.
}

// New creates and initializes a new Publisher instance.
func New(cfg config.NatsNotificationConfig) (*Publisher, error) {
	p := &Publisher{
		opts: getOptions(cfg),
	}

	// validate required fields.
	if err := p.validate(); err != nil {
		return nil, err
	}

	// connect to NATS server.
	nc, err := nats.Connect(p.opts.natsURL, p.getNatsOptions()...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNatsConnect, err)
	}
	p.nc = nc

	// create JetStream context.
	js, err := nc.JetStream(p.getJetStreamOptions()...)
	if err != nil {
		_ = p.nc.Drain() // cleanly close NATS on failure.
		return nil, fmt.Errorf("%w: %v", ErrGettingJetStream, err)
	}
	p.js = js

	return p, nil
}

// validate ensures required config fields are present.
func (p *Publisher) validate() error {
	switch {
	case p.opts.natsURL == "":
		return ErrNatsURLIsEmpty
	case p.opts.streamName == "":
		return ErrStreamNameIsEmpty
	case p.opts.subject == "":
		return ErrSubjectIsEmpty
	case p.opts.timeout == 0:
		return ErrTimeoutIsEmpty
	case p.opts.streamOption.name == "":
		return ErrNameInStreamOptionIsEmpty
	case p.opts.natsOption.name == "":
		return ErrNameInNatsOptionIsEmpty
	default:
		return nil
	}
}

// Start serializes and publishes the notification to JetStream with retries.
func (p *Publisher) Start(ctx context.Context, telegramID string, data notification.Notification) error {
	subject := p.opts.subject + telegramID

	// ensure the stream exists or create it.
	if err := p.ensureStream(); err != nil {
		return fmt.Errorf("%w: %v", ErrEnsureStreamFailed, err)
	}

	// serialize the message using msgpack.
	rawData, err := msgpack.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToMarshalNotification, err)
	}

	const maxRetries = 3
	var ackFuture nats.PubAckFuture

	// retry publish with exponential backoff.
	for i := 1; i <= maxRetries; i++ {
		const timeout = 100 // base backoff delay.

		// log high async pending threshold.
		if pending := p.js.PublishAsyncPending(); pending > p.opts.jetStreamOption.publishAsyncMaxPending-10 {
			log.Printf("[publisher][%s] high pending async messages: %d", telegramID, pending)
		}

		ackFuture, err = p.js.PublishAsync(subject, rawData)
		if err == nil {
			break
		}
		if i == maxRetries {
			return fmt.Errorf("failed to publish async after %d attempts: %w", i, err)
		}
		time.Sleep(time.Duration(i*timeout) * time.Millisecond)
	}

	// wait for acknowledgment with context timeout.
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(p.opts.timeout)*time.Second)
	defer cancel()

	select {
	case <-ctxTimeout.Done():
		return fmt.Errorf("%w [%s]", ErrPublishAsyncTimeoutForSubject, subject)
	case ack := <-ackFuture.Ok():
		if ack == nil {
			return ErrAckFutureOkReturnedNil
		}
		if ack.Stream != p.opts.streamName {
			return fmt.Errorf("%w: %s", ErrReceivedAckForUnexpectedStream, ack.Stream)
		}

		// check for publish error (non-blocking).
		select {
		case err := <-ackFuture.Err():
			if err != nil {
				return fmt.Errorf("%w: %v", ErrPublishAckFailed, err)
			}
		default:
			// no error to handle.
		}
		return nil
	}
}

// Close gracefully shuts down the publisher.
func (p *Publisher) Close() error {
	if p.nc != nil && !p.nc.IsClosed() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.opts.timeout)*time.Second)
		defer cancel()

		// wait for all async publishes to complete.
		if err := p.WaitForCompletion(ctx); err != nil {
			log.Printf("[Publisher] async publish wait failed: %v", err)
		}

		// cleanly drain NATS connection.
		if err := p.nc.Drain(); err != nil {
			log.Printf("[Publisher] drain failed: %v", err)
		}
	}
	return nil
}

// WaitForCompletion blocks until all async publish operations are acknowledged.
func (p *Publisher) WaitForCompletion(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		_ = p.js.PublishAsyncComplete()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ErrTimeoutWaitingForAsyncCompletion
	case <-done:
		return nil
	}
}

// ensureStream verifies that the stream exists; creates it if not.
func (p *Publisher) ensureStream() error {
	if _, err := p.js.StreamInfo(p.opts.streamName); err == nil {
		return nil // stream exists.
	}

	log.Printf("creating new stream [%s]...", p.opts.streamName)

	// attempt to add the stream.
	if _, err := p.js.AddStream(p.getStreamOptions()); err != nil {
		log.Printf("failed to add stream: %v", err)
		return err
	}

	log.Printf("✅ stream [%s] created successfully", p.opts.streamName)
	return nil
}

// getNatsOptions returns connection options for NATS.
func (p *Publisher) getNatsOptions() []nats.Option {
	opts := []nats.Option{
		nats.MaxReconnects(p.opts.natsOption.maxReconnects),
		nats.ReconnectWait(p.opts.natsOption.reconnectWait),
		nats.ReconnectJitter(p.opts.natsOption.reconnectJitter.jitter, p.opts.natsOption.reconnectJitter.jitterForTLS),
		nats.Name(p.opts.natsOption.name),
	}

	if p.opts.natsOption.errorHandler {
		opts = append(opts, nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("[Publisher] async error: %v", err)
		}))
	}

	return opts
}

// getJetStreamOptions returns JetStream-specific options.
func (p *Publisher) getJetStreamOptions() []nats.JSOpt {
	opts := []nats.JSOpt{
		nats.PublishAsyncMaxPending(p.opts.jetStreamOption.publishAsyncMaxPending),
	}

	if p.opts.jetStreamOption.publishAsyncErrHandler {
		opts = append(opts, nats.PublishAsyncErrHandler(func(_ nats.JetStream, _ *nats.Msg, err error) {
			log.Printf("[Publisher] async publish error: %v", err)
		}))
	}

	return opts
}

// getStreamOptions returns the stream configuration.
func (p *Publisher) getStreamOptions() *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:        p.opts.streamName,
		Subjects:    []string{p.opts.subject + "*"},
		Storage:     nats.StorageType(p.opts.streamOption.storage),
		Retention:   nats.RetentionPolicy(p.opts.streamOption.retention),
		MaxAge:      p.opts.streamOption.maxAge,
		MaxMsgs:     p.opts.streamOption.maxMsgs,
		MaxBytes:    p.opts.streamOption.maxBytes,
		Discard:     nats.DiscardPolicy(p.opts.streamOption.discard),
		Replicas:    p.opts.streamOption.replicas, // change to 3 if you need high availability. set to 3 for HA setup.
		AllowDirect: p.opts.streamOption.allowDirect,
		DenyDelete:  p.opts.streamOption.denyDelete,
		DenyPurge:   p.opts.streamOption.denyPurge,
	}
}
