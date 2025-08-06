package consumer

import (
	"time"

	"github.com/go-jedi/lingramm_backend/config"
)

type streamOption struct {
	maxMsgs   int64
	maxBytes  int64
	maxAge    time.Duration
	storage   int
	retention int
	discard   int
	name      string
}

type reconnectJitter struct {
	jitter       time.Duration
	jitterForTLS time.Duration
}

type natsOption struct {
	reconnectJitter      reconnectJitter
	reconnectWait        time.Duration
	timeout              time.Duration
	drainTimeout         time.Duration
	pingInterval         time.Duration
	maxReconnects        int
	maxPingsOutstanding  int
	closedHandler        bool
	reconnectHandler     bool
	disconnectErrHandler bool
	errorHandler         bool
	name                 string
}

type subscribeOption struct {
	ackWait       time.Duration
	maxAckPending int
	manualAck     bool
	deliverAll    bool
	durable       string
}

type options struct {
	streamOption    streamOption
	natsOption      natsOption
	subscribeOption subscribeOption
	natsURL         string
	streamName      string
	subject         string
	timeout         int
}

func getOptions(cfg config.NatsNotificationConfig) options {
	return options{
		streamOption: streamOption{
			maxMsgs:   cfg.Consumer.StreamOption.MaxMsgs,
			maxBytes:  cfg.Consumer.StreamOption.MaxBytes,
			maxAge:    time.Duration(cfg.Consumer.StreamOption.MaxAge) * time.Minute,
			storage:   cfg.Consumer.StreamOption.Storage,
			retention: cfg.Consumer.StreamOption.Retention,
			discard:   cfg.Consumer.StreamOption.Discard,
			name:      cfg.Consumer.StreamOption.Name,
		},
		natsOption: natsOption{
			reconnectJitter: reconnectJitter{
				jitter:       time.Duration(cfg.Consumer.NatsOption.ReconnectJitter.Jitter) * time.Millisecond,
				jitterForTLS: time.Duration(cfg.Consumer.NatsOption.ReconnectJitter.JitterForTLS) * time.Millisecond,
			},
			reconnectWait:        time.Duration(cfg.Consumer.NatsOption.ReconnectWait) * time.Second,
			timeout:              time.Duration(cfg.Consumer.NatsOption.Timeout) * time.Second,
			drainTimeout:         time.Duration(cfg.Consumer.NatsOption.DrainTimeout) * time.Second,
			pingInterval:         time.Duration(cfg.Consumer.NatsOption.PingInterval) * time.Minute,
			maxReconnects:        cfg.Consumer.NatsOption.MaxReconnects,
			maxPingsOutstanding:  cfg.Consumer.NatsOption.MaxPingsOutstanding,
			closedHandler:        cfg.Consumer.NatsOption.ClosedHandler,
			reconnectHandler:     cfg.Consumer.NatsOption.ReconnectHandler,
			disconnectErrHandler: cfg.Consumer.NatsOption.DisconnectErrHandler,
			errorHandler:         cfg.Consumer.NatsOption.ErrorHandler,
			name:                 cfg.Consumer.NatsOption.Name,
		},
		subscribeOption: subscribeOption{
			ackWait:       time.Duration(cfg.Consumer.SubscribeOption.AckWait) * time.Second,
			maxAckPending: cfg.Consumer.SubscribeOption.MaxAckPending,
			manualAck:     cfg.Consumer.SubscribeOption.ManualAck,
			deliverAll:    cfg.Consumer.SubscribeOption.DeliverAll,
			durable:       cfg.Consumer.SubscribeOption.Durable,
		},
		natsURL:    cfg.NatsURL,
		streamName: cfg.Consumer.StreamName,
		subject:    cfg.Consumer.Subject,
		timeout:    cfg.Consumer.Timeout,
	}
}
