package broker

import (
	"context"
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	dsn          string
	exchangeName string
	exchangeType string
	routingKey   string
	name         string
	logger       logger.ILogger
}

func NewProducer(
	cfg *config.Config,
	logger logger.ILogger,
) *Producer {
	return &Producer{
		logger:       logger,
		dsn:          cfg.Publisher.Dsn,
		name:         cfg.Publisher.QueueName,
		exchangeName: cfg.Publisher.ExchangeName,
		exchangeType: cfg.Publisher.ExchangeType,
		routingKey:   cfg.Publisher.RoutingKey,
	}
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	cfg := amqp.Config{Properties: amqp.NewConnectionProperties()}
	cfg.Properties.SetClientConnectionName("sample-producer")
	p.logger.Infof("dialing %q", p.dsn)
	connection, err := amqp.DialConfig(p.dsn, cfg)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer connection.Close()

	p.logger.Info("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	p.logger.Infof("got Channel, declaring %q Exchange (%q)", p.exchangeType, p.exchangeName)
	if err := channel.ExchangeDeclare(
		p.exchangeName, // queueName
		p.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %w", err)
	}

	p.logger.Info("declared Exchange, publishing messages")

	p.logger.Infof("publishing %dB body (%q)", len(body), body)

	err = channel.PublishWithContext(ctx,
		p.exchangeName, // publish to an exchangeName
		p.routingKey,   // routing to 0 or more queues
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	)

	if err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}

	p.logger.Infof("published %dB OK", len(body))

	return nil
}
