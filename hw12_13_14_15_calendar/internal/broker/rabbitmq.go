package broker

import (
	"context"
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	RabbitMQ struct {
		conn    *amqp.Connection
		channel *amqp.Channel
		dsn     string
		queue   string
		logger  logger.ILogger
	}
)

func New(cfg *config.Config, logger logger.ILogger) *RabbitMQ {
	return &RabbitMQ{
		dsn:    cfg.Broker.Dsn,
		queue:  cfg.Broker.Name,
		logger: logger,
	}
}

func (b *RabbitMQ) Connect(ctx context.Context) error {
	var scopeErr error

	b.logger.Info("broker is connecting...")

	b.conn, scopeErr = amqp.Dial(b.dsn)

	if scopeErr != nil {
		return fmt.Errorf("connect to broker failed: %w", scopeErr)
	}

	b.channel, scopeErr = b.conn.Channel()

	if scopeErr != nil {
		return fmt.Errorf("make channel failed: %w", scopeErr)
	}

	_, scopeErr = b.channel.QueueDeclare(
		b.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if scopeErr != nil {
		return fmt.Errorf("queue declade failed: %w", scopeErr)
	}

	b.logger.Info("broker is connect")

	return nil
}

func (b *RabbitMQ) Close() error {
	b.logger.Info("broker is closing")

	err := b.channel.Close()
	if err != nil {
		return fmt.Errorf("channel close failed: %w", err)
	}
	err = b.conn.Close()
	if err != nil {
		return fmt.Errorf("broker connection close failed: %w", err)
	}

	b.logger.Info("broker is close")
	return err
}
