package broker

import (
	"errors"
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connector interface {
	Connect() error
	Reconnect() error
	BindQueue() error
}

type RMQConfig struct {
	Tag, ConnectionName string
}

// Connection is the connection created.
type Connection struct {
	logger         logger.Logger
	connectionName string
	dsn            string
	conn           *amqp.Connection
	channel        *amqp.Channel
	exchange       string
	exchangeType   string
	queue          string
	tag            string
	bindingKey     string
	err            chan error
}

// NewConnection returns the new connection object.
func NewConnection(
	rmqConfig RMQConfig,
	logger logger.Logger,
	cfg *config.Config,
) *Connection {
	return &Connection{
		logger:         logger,
		exchange:       cfg.RMQ.ExchangeName,
		exchangeType:   cfg.RMQ.ExchangeType,
		dsn:            cfg.RMQ.Dsn,
		queue:          cfg.RMQ.QueueName,
		connectionName: rmqConfig.ConnectionName,
		tag:            rmqConfig.Tag,
		bindingKey:     cfg.RMQ.BindingKey,
		err:            make(chan error),
	}
}

func (c *Connection) Connect() error {
	var err error

	cfg := amqp.Config{Properties: amqp.NewConnectionProperties()}
	cfg.Properties.SetClientConnectionName(c.connectionName)
	c.logger.Infof("dialing %q", c.dsn)
	c.conn, err = amqp.DialConfig(c.dsn, cfg)
	if err != nil {
		return fmt.Errorf("error in creating rabbitmq connection with %s : %w", c.dsn, err)
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error))
		c.err <- errors.New("connection Closed")
	}()
	c.logger.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	c.logger.Infof("got Channel, declaring %q Exchange (%q)", c.exchangeType, c.exchange)
	if err := c.channel.ExchangeDeclare(
		c.exchange,     // name
		c.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("error in Exchange Declare: %w", err)
	}
	c.logger.Info("declared Exchange")
	return nil
}

func (c *Connection) BindQueue() error {
	queue, err := c.channel.QueueDeclare(
		c.queue, // queueName of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %w", err)
	}

	c.logger.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, c.bindingKey)

	if err = c.channel.QueueBind(
		queue.Name,   // queueName of the queue
		c.bindingKey, // bindingKey
		c.exchange,   // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("queue Bind: %w", err)
	}

	c.logger.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)

	return nil
}

// Reconnect reconnects the connection.
func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}
