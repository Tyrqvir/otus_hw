package broker

import (
	"context"
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, body []byte) error
}

type Producer struct {
	exchangeName string
	routingKey   string
	logger       logger.Logger
	connection   *Connection
}

func NewProducer(
	cfg *config.Config,
	logger logger.Logger,
	connection *Connection,
) *Producer {
	return &Producer{
		logger:       logger,
		exchangeName: cfg.RMQ.ExchangeName,
		routingKey:   cfg.RMQ.BindingKey,
		connection:   connection,
	}
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	select {
	case err := <-p.connection.err:
		if err != nil {
			p.connection.Reconnect()
		}
	default:
	}

	p.logger.Infof("publishing %dB body (%q)", len(body), body)

	err := p.connection.channel.PublishWithContext(ctx,
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
