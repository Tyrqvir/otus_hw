package scheduler

import (
	"context"
	"os"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

type (
	IQueue interface {
		Connect(ctx context.Context) error
		Close() error
	}

	Scheduler struct {
		logger logger.ILogger
		queue  IQueue
	}
)

func New(
	cfg *config.Config,
	logger logger.ILogger,
	queue IQueue,
) *Scheduler {
	return &Scheduler{
		logger: logger,
		queue:  queue,
	}
}

func (s *Scheduler) Run(ctx context.Context, stop context.CancelFunc) error {
	s.logger.Info("start scheduler...")

	go func() {
		if err := s.queue.Connect(ctx); err != nil {
			s.logger.Error("connect to broker: " + err.Error())
			stop()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		if err := s.queue.Close(); err != nil {
			s.logger.Error("close broker: " + err.Error())
			stop()
			os.Exit(1)
		}
	}()

	return nil
}
