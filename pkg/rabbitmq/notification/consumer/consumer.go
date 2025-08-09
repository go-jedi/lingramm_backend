package consumer

import (
	"context"
	"errors"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	defaultTimeoutCheckConnect   = 20
	defaultTimeoutPublishMessage = 5
)

var (
	ErrURLRabbitmqNotFound  = errors.New("url amqp not found")
	ErrExchangeNameNotFound = errors.New("exchange name not found")
	ErrExchangeKindNotFound = errors.New("exchange kind not found")
)

// IConsumer defines the interface for the consumer notification.
//
//go:generate mockery --name=IConsumer --output=mocks --case=underscore
type IConsumer interface {
	Execute(ctx context.Context, telegramID string, fn func(msg notification.SendNotificationDTO) error) error
	Close() error
}

type Consumer struct {
	conn                  *amqp.Connection
	timeoutCheckConnect   time.Duration // timeout to check connection.
	timeoutPublishMessage time.Duration // timeout publish message.
	opts                  Options
}

func New(cfg config.ConsumerConfig) (*Consumer, error) {
	c := &Consumer{
		opts: getOptions(cfg),
	}

	if err := c.init(); err != nil {
		return nil, err
	}

	// connection to rabbitmq.
	conn, err := amqp.DialConfig(c.opts.URL, c.getAmqpConfig())
	if err != nil {
		return nil, err
	}
	c.conn = conn

	// call method to check connection rabbitmq in background.
	go c.checkConnection()

	return c, nil
}

func (c *Consumer) init() error {
	if c.opts.URL == "" {
		return ErrURLRabbitmqNotFound
	}
	if c.opts.Exchange.Name == "" {
		return ErrExchangeNameNotFound
	}
	if c.opts.Exchange.Kind == "" {
		return ErrExchangeKindNotFound
	}

	if c.timeoutCheckConnect == 0 {
		c.timeoutCheckConnect = defaultTimeoutCheckConnect
	}
	if c.timeoutPublishMessage == 0 {
		c.timeoutPublishMessage = defaultTimeoutPublishMessage
	}

	return nil
}

func (c *Consumer) Execute(_ context.Context, telegramID string, fn func(msg notification.SendNotificationDTO) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare exchange. if not exists created.
	if err := c.exchangeDeclare(ch); err != nil {
		return err
	}

	// declare queue.
	q, err := c.queueDeclare(ch)
	if err != nil {
		return err
	}

	// queue bind.
	if err := c.queueBind(ch, q, telegramID); err != nil {
		return err
	}

	// consume to read messages.
	if err := c.consume(ch, q, fn); err != nil {
		return err
	}

	return nil
}

// Close connection rabbitmq.
func (c *Consumer) Close() error {
	return c.conn.Close()
}

// consume immediately starts delivering queued messages.
func (c *Consumer) consume(channel *amqp.Channel, queue amqp.Queue, fn func(msg notification.SendNotificationDTO) error) error {
	delivery, err := channel.Consume(
		queue.Name,
		c.opts.Consume.Consumer,
		c.opts.Consume.AutoAck,
		c.opts.Consume.Exclusive,
		c.opts.Consume.NoLocal,
		c.opts.Consume.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	chExit := make(chan struct{})

	go c.messageHandler(delivery, chExit, fn)

	<-chExit

	return nil
}

// messageHandler message handler.
func (c *Consumer) messageHandler(delivery <-chan amqp.Delivery, chExit chan<- struct{}, fn func(msg notification.SendNotificationDTO) error) {
	const contentTypeMsgpack = "application/msgpack"

	for d := range delivery {
		if d.ContentType != contentTypeMsgpack {
			log.Printf("[consumer]: unexpected content-type %q", d.ContentType)
			if err := d.Nack(false, false); err != nil {
				log.Printf("[consumer]: failed to nack message: %v", err)
			}
			continue
		}

		var result notification.SendNotificationDTO
		if err := msgpack.Unmarshal(d.Body, &result); err != nil {
			log.Printf("[consumer]: failed to unmarshal message: %v", err)
			if err := d.Nack(false, false); err != nil {
				log.Printf("[consumer]: failed to nack message: %v", err)
			}
			continue
		}

		if err := fn(result); err != nil {
			log.Printf("[consumer]: failed to message transmission: %v", err)
			if err := d.Nack(false, false); err != nil {
				log.Printf("[consumer]: failed to nack message: %v", err)
			}
			continue
		}

		if err := d.Ack(false); err != nil {
			log.Printf("[consumer]: failed to ack message: %v", err)
		}
	}

	close(chExit)
}

// exchangeDeclare declare exchange. if not exists created exchange.
func (c *Consumer) exchangeDeclare(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		c.opts.Exchange.Name,
		c.opts.Exchange.Kind,
		c.opts.Exchange.Durable,
		c.opts.Exchange.AutoDelete,
		c.opts.Exchange.Internal,
		c.opts.Exchange.NoWait,
		nil,
	)
}

// queueDeclare queue declare.
func (c *Consumer) queueDeclare(channel *amqp.Channel) (amqp.Queue, error) {
	const (
		xExpires    = "x-expires"
		xMessageTTL = "x-message-ttl"
	)

	args := amqp.Table{}

	if c.opts.Queue.Args.XExpires > 0 {
		args[xExpires] = c.opts.Queue.Args.XExpires
	}
	if c.opts.Queue.Args.XMessageTTL > 0 {
		args[xMessageTTL] = c.opts.Queue.Args.XMessageTTL
	}

	return channel.QueueDeclare(
		c.opts.Queue.Name,
		c.opts.Queue.Durable,
		c.opts.Queue.AutoDelete,
		c.opts.Queue.Exclusive,
		c.opts.Queue.NoWait,
		args,
	)
}

// queueBind queue bind by telegram id.
func (c *Consumer) queueBind(channel *amqp.Channel, queue amqp.Queue, telegramID string) error {
	return channel.QueueBind(
		queue.Name,
		telegramID,
		c.opts.QueueBind.Exchange,
		c.opts.QueueBind.NoWait,
		nil,
	)
}

// reconnection rabbitmq.
func (c *Consumer) reconnection() error {
	conn, err := amqp.DialConfig(c.opts.URL, c.getAmqpConfig())
	if err != nil {
		return err
	}
	c.conn = conn

	log.Println("reconnected consumer to rabbitmq")

	return nil
}

// checkConnection check connection rabbitmq and recovery.
func (c *Consumer) checkConnection() {
	for {
		if c.conn.IsClosed() {
			if err := c.reconnection(); err != nil {
				log.Printf("reconnection consumer error: %v", err)
			}
		}
		time.Sleep(c.timeoutCheckConnect)
	}
}

// getAmqpConfig get amqp config for connection to rabbitmq.
func (c *Consumer) getAmqpConfig() amqp.Config {
	const connectionName = "connection_name"

	var (
		table = amqp.Table{}
		sasl  []amqp.Authentication
	)

	if c.opts.Amqp.SASL.IsPlainAuth {
		sasl = append(sasl, &amqp.PlainAuth{
			Username: c.opts.Amqp.SASL.PlainAuth.Username,
			Password: c.opts.Amqp.SASL.PlainAuth.Password,
		})
	}

	if len(c.opts.Amqp.Properties.ConnectionName) > 0 {
		table[connectionName] = c.opts.Amqp.Properties.ConnectionName
	}

	opts := amqp.Config{
		SASL:       sasl,                   // SASL не обязательно задавать, если логин/пароль есть в URL
		ChannelMax: c.opts.Amqp.ChannelMax, // ограниченный верхний предел каналов на коннект (0 = без ограничения)
		FrameSize:  c.opts.Amqp.FrameSize,  // 0 => дефолт (131072). Поднимай только если гоняешь крупные сообщения (МБ)
		Heartbeat:  c.opts.Amqp.HeartBeat,  // быстрый детект обрывов.
		Properties: table,
		Locale:     c.opts.Amqp.Locale,
	}

	if c.opts.Amqp.IsDial {
		dialer := &net.Dialer{
			Timeout:   c.opts.Amqp.Dialer.Timeout,   // быстро падаем при недоступности
			KeepAlive: c.opts.Amqp.Dialer.KeepAlive, // удерживаем соединение живым
			Control: func(_ string, _ string, c syscall.RawConn) error {
				var opErr error
				if err := c.Control(func(fd uintptr) {
					// Отключаем Nagle (отправка сразу без накопления)
					opErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
					if opErr != nil {
						return
					}
					// Ускоряем отправку ACK
					opErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1)
				}); err != nil {
					return err
				}
				return opErr
			},
		}
		opts.Dial = dialer.Dial
	}

	return opts
}
