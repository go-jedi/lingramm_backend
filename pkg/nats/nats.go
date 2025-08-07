package nats

import (
	"github.com/go-jedi/lingramm_backend/config"
	notificationconsumer "github.com/go-jedi/lingramm_backend/pkg/nats/notification/consumer"
	notificationpublisher "github.com/go-jedi/lingramm_backend/pkg/nats/notification/publisher"
)

type Nats struct {
	Notification struct {
		Consumer  notificationconsumer.IConsumer
		Publisher notificationpublisher.IPublisher
	}
}

func New(cfg config.NatsConfig) (*Nats, error) {
	n := &Nats{}

	nc, err := notificationconsumer.New(cfg.Notification)
	if err != nil {
		return nil, err
	}
	n.Notification.Consumer = nc

	np, err := notificationpublisher.New(cfg.Notification)
	if err != nil {
		return nil, err
	}
	n.Notification.Publisher = np

	return n, nil
}
