package broker

import (
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	logger       logger.Logger
	dsn          string
	exchangeName string
	exchangeType string
	bindingKey   string
	queueName    string
	tag          string
}

func NewConsumer(
	cfg *config.Config,
	logger logger.Logger,
) *Consumer {
	return &Consumer{
		logger:       logger,
		dsn:          cfg.Consumer.Dsn,
		queueName:    cfg.Consumer.QueueName,
		exchangeName: cfg.Consumer.ExchangeName,
		exchangeType: cfg.Consumer.ExchangeType,
		bindingKey:   cfg.Consumer.BindingKey,
		tag:          "simple-consumer",
	}
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	var err error
	cfg := amqp.Config{Properties: amqp.NewConnectionProperties()}
	cfg.Properties.SetClientConnectionName("sample-consumer")
	c.logger.Infof("dialing %q", c.dsn)
	c.conn, err = amqp.DialConfig(c.dsn, cfg)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c.logger.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}

	c.logger.Infof("got Channel, declaring Exchange (%q)", c.exchangeName)
	if err = c.channel.ExchangeDeclare(
		c.exchangeName, // queueName of the exchangeName
		c.exchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %w", err)
	}

	c.logger.Infof("declared Exchange, declaring Queue %q", c.queueName)
	queue, err := c.channel.QueueDeclare(
		c.queueName, // queueName of the queue
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %w", err)
	}

	c.logger.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, c.bindingKey)

	if err = c.channel.QueueBind(
		queue.Name,     // queueName of the queue
		c.bindingKey,   // bindingKey
		c.exchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %w", err)
	}

	c.logger.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // queueName
		c.tag,      // consumerTag,
		true,       // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %w", err)
	}

	return deliveries, nil
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	defer c.logger.Info("AMQP shutdown OK")

	return nil
}
