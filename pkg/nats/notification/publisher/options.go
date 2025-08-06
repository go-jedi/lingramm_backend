package publisher

import (
	"time"

	"github.com/go-jedi/lingramm_backend/config"
)

type streamOption struct {
	maxMsgs     int64
	maxBytes    int64
	maxAge      time.Duration
	storage     int
	retention   int
	discard     int
	replicas    int
	allowDirect bool
	denyDelete  bool
	denyPurge   bool
	name        string
}

type reconnectJitter struct {
	jitter       time.Duration
	jitterForTLS time.Duration
}

type natsOption struct {
	reconnectJitter reconnectJitter
	reconnectWait   time.Duration
	maxReconnects   int
	errorHandler    bool
	name            string
}

type jetStreamOption struct {
	publishAsyncMaxPending int
	publishAsyncErrHandler bool
}

type options struct {
	streamOption    streamOption
	natsOption      natsOption
	jetStreamOption jetStreamOption
	timeout         int
	natsURL         string
	streamName      string
	subject         string
}

func getOptions(cfg config.NatsNotificationConfig) options {
	return options{
		streamOption: streamOption{
			maxMsgs:     cfg.Publisher.StreamOption.MaxMsgs,
			maxBytes:    cfg.Publisher.StreamOption.MaxBytes,
			maxAge:      time.Duration(cfg.Publisher.StreamOption.MaxAge) * time.Minute,
			storage:     cfg.Publisher.StreamOption.Storage,
			retention:   cfg.Publisher.StreamOption.Retention,
			discard:     cfg.Publisher.StreamOption.Discard,
			replicas:    cfg.Publisher.StreamOption.Replicas,
			allowDirect: cfg.Publisher.StreamOption.AllowDirect,
			denyDelete:  cfg.Publisher.StreamOption.DenyDelete,
			denyPurge:   cfg.Publisher.StreamOption.DenyPurge,
			name:        cfg.Publisher.StreamOption.Name,
		},
		natsOption: natsOption{
			reconnectJitter: reconnectJitter{
				jitter:       time.Duration(cfg.Publisher.NatsOption.ReconnectJitter.Jitter) * time.Millisecond,
				jitterForTLS: time.Duration(cfg.Publisher.NatsOption.ReconnectJitter.JitterForTLS) * time.Millisecond,
			},
			reconnectWait: time.Duration(cfg.Publisher.NatsOption.ReconnectWait) * time.Second,
			maxReconnects: cfg.Publisher.NatsOption.MaxReconnects,
			errorHandler:  cfg.Publisher.NatsOption.ErrorHandler,
			name:          cfg.Publisher.NatsOption.Name,
		},
		jetStreamOption: jetStreamOption{
			publishAsyncMaxPending: cfg.Publisher.JetStreamOption.PublishAsyncMaxPending,
			publishAsyncErrHandler: cfg.Publisher.JetStreamOption.PublishAsyncErrHandler,
		},
		timeout:    cfg.Publisher.Timeout,
		natsURL:    cfg.NatsURL,
		streamName: cfg.Publisher.StreamName,
		subject:    cfg.Publisher.Subject,
	}
}
