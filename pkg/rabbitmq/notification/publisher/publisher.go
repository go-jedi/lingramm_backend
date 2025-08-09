package publisher

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
	ErrURLRabbitmqNotFound           = errors.New("url amqp not found")
	ErrExchangeNameNotFound          = errors.New("exchange name not found")
	ErrExchangeKindNotFound          = errors.New("exchange kind not found")
	ErrPublishExchangeNotFound       = errors.New("publish exchange not found")
	ErrPublishingContentTypeNotFound = errors.New("publishing content type not found")
	ErrPublishingTypeNotFound        = errors.New("publishing type not found")
)

// IPublisher defines the interface for the publisher notification.
//
//go:generate mockery --name=IPublisher --output=mocks --case=underscore
type IPublisher interface {
	Execute(ctx context.Context, telegramID string, data notification.SendNotificationDTO) error
	Close() error
}

type Publisher struct {
	conn                  *amqp.Connection
	timeoutCheckConnect   time.Duration // timeout to check connection.
	timeoutPublishMessage time.Duration // timeout publish message.
	opts                  Options
}

func New(cfg config.PublisherConfig) (*Publisher, error) {
	p := &Publisher{
		opts: getOptions(cfg),
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	// connection to rabbitmq.
	conn, err := amqp.DialConfig(p.opts.URL, p.getAmqpConfig())
	if err != nil {
		return nil, err
	}
	p.conn = conn

	// call method to check connection rabbitmq in background.
	go p.checkConnection()

	return p, nil
}

func (p *Publisher) init() error {
	if p.opts.URL == "" {
		return ErrURLRabbitmqNotFound
	}
	if p.opts.Exchange.Name == "" {
		return ErrExchangeNameNotFound
	}
	if p.opts.Exchange.Kind == "" {
		return ErrExchangeKindNotFound
	}
	if p.opts.Publish.Exchange == "" {
		return ErrPublishExchangeNotFound
	}
	if p.opts.Publish.Publishing.ContentType == "" {
		return ErrPublishingContentTypeNotFound
	}
	if p.opts.Publish.Publishing.Type == "" {
		return ErrPublishingTypeNotFound
	}

	if p.timeoutCheckConnect == 0 {
		p.timeoutCheckConnect = defaultTimeoutCheckConnect
	}
	if p.timeoutPublishMessage == 0 {
		p.timeoutPublishMessage = defaultTimeoutPublishMessage
	}

	return nil
}

func (p *Publisher) Execute(ctx context.Context, telegramID string, data notification.SendNotificationDTO) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare exchange. if not exists created.
	if err := p.exchangeDeclare(ch); err != nil {
		return err
	}

	// marshal data to msgpack bytes.
	rawData, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	// publish message with context.
	if err := p.publishWithContext(ctx, ch, telegramID, rawData); err != nil {
		return err
	}

	return nil
}

// Close connection rabbitmq.
func (p *Publisher) Close() error {
	return p.conn.Close()
}

// exchangeDeclare declare exchange. if not exists created exchange.
func (p *Publisher) exchangeDeclare(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		p.opts.Exchange.Name,
		p.opts.Exchange.Kind,
		p.opts.Exchange.Durable,
		p.opts.Exchange.AutoDelete,
		p.opts.Exchange.Internal,
		p.opts.Exchange.NoWait,
		nil,
	)
}

// publishWithContext publish message with context.
func (p *Publisher) publishWithContext(ctx context.Context, channel *amqp.Channel, telegramID string, body []byte) error {
	const timestamp = "now"

	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeoutPublishMessage)
	defer cancel()

	msg := amqp.Publishing{
		ContentType: p.opts.Publish.Publishing.ContentType,
		Type:        p.opts.Publish.Publishing.Type,
		Body:        body,
	}

	if p.opts.Publish.Publishing.Timestamp == timestamp {
		msg.Timestamp = time.Now()
	}

	return channel.PublishWithContext(
		ctxTimeout,
		p.opts.Publish.Exchange,
		telegramID,
		p.opts.Publish.Mandatory,
		p.opts.Publish.Immediate,
		msg,
	)
}

// reconnection rabbitmq.
func (p *Publisher) reconnection() error {
	conn, err := amqp.DialConfig(p.opts.URL, p.getAmqpConfig())
	if err != nil {
		return err
	}
	p.conn = conn

	log.Println("reconnected publisher to rabbitmq")

	return nil
}

// checkConnection check connection rabbitmq and recovery.
func (p *Publisher) checkConnection() {
	for {
		if p.conn.IsClosed() {
			if err := p.reconnection(); err != nil {
				log.Printf("reconnection publisher error: %v", err)
			}
		}
		time.Sleep(p.timeoutCheckConnect)
	}
}

// getAmqpConfig get amqp config for connection to rabbitmq.
func (p *Publisher) getAmqpConfig() amqp.Config {
	const connectionName = "connection_name"

	var (
		table = amqp.Table{}
		sasl  []amqp.Authentication
	)

	if p.opts.Amqp.SASL.IsPlainAuth {
		sasl = append(sasl, &amqp.PlainAuth{
			Username: p.opts.Amqp.SASL.PlainAuth.Username,
			Password: p.opts.Amqp.SASL.PlainAuth.Password,
		})
	}

	if len(p.opts.Amqp.Properties.ConnectionName) > 0 {
		table[connectionName] = p.opts.Amqp.Properties.ConnectionName
	}

	opts := amqp.Config{
		SASL:       sasl,                   // SASL не обязательно задавать, если логин/пароль есть в URL
		ChannelMax: p.opts.Amqp.ChannelMax, // ограниченный верхний предел каналов на коннект (0 = без ограничения)
		FrameSize:  p.opts.Amqp.FrameSize,  // 0 => дефолт (131072). Поднимай только если гоняешь крупные сообщения (МБ)
		Heartbeat:  p.opts.Amqp.HeartBeat,  // быстрый детект обрывов.
		Properties: table,
		Locale:     p.opts.Amqp.Locale,
	}

	if p.opts.Amqp.IsDial {
		dialer := &net.Dialer{
			Timeout:   p.opts.Amqp.Dialer.Timeout,   // быстро падаем при недоступности
			KeepAlive: p.opts.Amqp.Dialer.KeepAlive, // удерживаем соединение живым
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
