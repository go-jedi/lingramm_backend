package rabbitmq

import (
	"github.com/go-jedi/lingramm_backend/config"
	notificationconsumer "github.com/go-jedi/lingramm_backend/pkg/rabbitmq/notification/consumer"
	notificationpublisher "github.com/go-jedi/lingramm_backend/pkg/rabbitmq/notification/publisher"
)

type RabbitMQ struct {
	Notification struct {
		Consumer  notificationconsumer.IConsumer
		Publisher notificationpublisher.IPublisher
	}
}

func New(cfg config.RabbitMQConfig) (*RabbitMQ, error) {
	r := &RabbitMQ{}

	nc, err := notificationconsumer.New(cfg.Notification.Consumer)
	if err != nil {
		return nil, err
	}
	r.Notification.Consumer = nc

	np, err := notificationpublisher.New(cfg.Notification.Publisher)
	if err != nil {
		return nil, err
	}
	r.Notification.Publisher = np

	return r, nil
}
