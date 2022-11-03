package broker

import (
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *Connection
	logger     logger.Logger
	queueName  string
	tag        string
}

func NewConsumer(
	cfg *config.Config,
	logger logger.Logger,
	connection *Connection,
) *Consumer {
	return &Consumer{
		logger:     logger,
		queueName:  cfg.RMQ.QueueName,
		tag:        cfg.Consumer.Tag,
		connection: connection,
	}
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	var err error

	deliveries, err := c.connection.channel.Consume(
		c.queueName, // queueName
		c.tag,       // consumerTag,
		true,        // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %w", err)
	}

	return deliveries, nil
}

func (c *Consumer) Shutdown() error {
	if err := c.connection.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := c.connection.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	defer c.logger.Info("AMQP shutdown OK")

	return nil
}
