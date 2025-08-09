package publisher

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

type Publishing struct {
	ContentType string
	Timestamp   string
	Type        string
}

type Publish struct {
	Exchange   string
	Mandatory  bool
	Immediate  bool
	Publishing Publishing
}

type Options struct {
	URL                   string
	TimeoutCheckConnect   time.Duration
	TimeoutPublishMessage time.Duration
	Amqp                  Amqp
	Exchange              Exchange
	Publish               Publish
}

func getOptions(cfg config.PublisherConfig) Options {
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
		Publish: Publish{
			Exchange:  cfg.Publish.Exchange,
			Mandatory: cfg.Publish.Mandatory,
			Immediate: cfg.Publish.Immediate,
			Publishing: Publishing{
				ContentType: cfg.Publish.Publishing.ContentType,
				Timestamp:   cfg.Publish.Publishing.Timestamp,
				Type:        cfg.Publish.Publishing.Type,
			},
		},
	}
}
