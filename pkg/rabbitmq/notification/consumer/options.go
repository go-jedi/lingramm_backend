package consumer

import (
	"time"

	"github.com/go-jedi/lingramm_backend/config"
)

type Dialer struct {
	Timeout   time.Duration
	KeepAlive time.Duration
}

type PlainAuth struct {
	Username string
	Password string
}

type SASL struct {
	IsPlainAuth bool
	PlainAuth   PlainAuth
}

type Properties struct {
	ConnectionName string
}

type Amqp struct {
	Dialer     Dialer
	SASL       SASL
	VHost      string
	ChannelMax uint16
	FrameSize  int
	HeartBeat  time.Duration
	Properties Properties
	Locale     string
	IsDial     bool
}

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type Args struct {
	XExpires    int32
	XMessageTTL int32
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       Args
}

type QueueBind struct {
	Exchange string
	NoWait   bool
}

type Consume struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type Options struct {
	URL                   string
	TimeoutCheckConnect   time.Duration
	TimeoutPublishMessage time.Duration
	Amqp                  Amqp
	Exchange              Exchange
	Queue                 Queue
	QueueBind             QueueBind
	Consume               Consume
}

func getOptions(cfg config.ConsumerConfig) Options {
	return Options{
		URL:                   cfg.URL,
		TimeoutCheckConnect:   time.Duration(cfg.TimeoutCheckConnect) * time.Second,
		TimeoutPublishMessage: time.Duration(cfg.TimeoutPublishMessage) * time.Second,
		Amqp: Amqp{
			Dialer: Dialer{
				Timeout:   time.Duration(cfg.Amqp.Dialer.Timeout) * time.Second,
				KeepAlive: time.Duration(cfg.Amqp.Dialer.KeepAlive) * time.Second,
			},
			SASL: SASL{
				IsPlainAuth: cfg.Amqp.SASL.IsPlainAuth,
				PlainAuth: PlainAuth{
					Username: cfg.Amqp.SASL.PlainAuth.Username,
					Password: cfg.Amqp.SASL.PlainAuth.Password,
				},
			},
			VHost:      cfg.Amqp.VHost,
			ChannelMax: cfg.Amqp.ChannelMax,
			FrameSize:  cfg.Amqp.FrameSize,
			HeartBeat:  time.Duration(cfg.Amqp.HeartBeat) * time.Second,
			Properties: Properties{
				ConnectionName: cfg.Amqp.Properties.ConnectionName,
			},
			Locale: cfg.Amqp.Locale,
			IsDial: cfg.Amqp.IsDial,
		},
		Exchange: Exchange{
			Name:       cfg.Exchange.Name,
			Kind:       cfg.Exchange.Kind,
			Durable:    cfg.Exchange.Durable,
			AutoDelete: cfg.Exchange.AutoDelete,
			Internal:   cfg.Exchange.Internal,
			NoWait:     cfg.Exchange.NoWait,
		},
		Queue: Queue{
			Name:       cfg.Queue.Name,
			Durable:    cfg.Queue.Durable,
			AutoDelete: cfg.Queue.AutoDelete,
			Exclusive:  cfg.Queue.Exclusive,
			NoWait:     cfg.Queue.NoWait,
			Args: Args{
				XExpires:    cfg.Queue.Args.XExpires,
				XMessageTTL: cfg.Queue.Args.XMessageTTL,
			},
		},
		QueueBind: QueueBind{
			Exchange: cfg.QueueBind.Exchange,
			NoWait:   cfg.QueueBind.NoWait,
		},
		Consume: Consume{
			Consumer:  cfg.Consume.Consumer,
			AutoAck:   cfg.Consume.AutoAck,
			Exclusive: cfg.Consume.Exclusive,
			NoLocal:   cfg.Consume.NoLocal,
			NoWait:    cfg.Consume.NoWait,
		},
	}
}
