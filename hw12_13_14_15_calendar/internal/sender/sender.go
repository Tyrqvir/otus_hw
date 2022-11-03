package sender

import (
	"context"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Consumer interface {
		Consume() (<-chan amqp.Delivery, error)
		Shutdown() error
	}

	Sender struct {
		logger     logger.Logger
		consumer   *broker.Consumer
		connection *broker.Connection
	}
)

func New(
	logger logger.Logger,
	consumer *broker.Consumer,
	connection *broker.Connection,
) *Sender {
	return &Sender{
		logger:     logger,
		consumer:   consumer,
		connection: connection,
	}
}

func (s *Sender) Handle(ctx context.Context) {
	s.logger.Info("sender starting ...")

	if err := s.connection.Connect(); err != nil {
		s.logger.Error("connect:" + err.Error())
	}

	if err := s.connection.BindQueue(); err != nil {
		s.logger.Error("bind queue:" + err.Error())
	}

	s.logger.Info("consuming ...")

	msgs, err := s.consumer.Consume()
	if err != nil {
		s.logger.Errorf("can't consume: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var notice model.Notice
			if err := json.Unmarshal(msg.Body, &notice); err != nil {
				s.logger.Error("can't unmarshal:", "err", err.Error())
				continue
			}

			s.logger.Info("success received notice:", "Id", notice.ID, "title", notice.Title)
		}
	}
}

func (s *Sender) Shutdown() error {
	err := s.consumer.Shutdown()
	if err != nil {
		s.logger.Error("can't shutdown:", "err", err.Error())
	}

	return nil
}
