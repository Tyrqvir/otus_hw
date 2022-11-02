package scheduler

import (
	"context"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/goccy/go-json"
)

type (
	Publisher interface {
		Publish(ctx context.Context, body []byte) error
	}

	Scheduler struct {
		logger     logger.Logger
		publisher  Publisher
		repository repository.EventRepository
		interval   time.Duration
	}
)

func New(
	cfg *config.Config,
	logger logger.Logger,
	publisher Publisher,
	repository repository.EventRepository,
) *Scheduler {
	return &Scheduler{
		logger:     logger,
		publisher:  publisher,
		repository: repository,
		interval:   cfg.Schedule.Interval,
	}
}

func (s *Scheduler) Handle(ctx context.Context) {
	s.logger.Info("scheduler starting ...")

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		go func() {
			date := time.Now()
			notifications, err := s.repository.NoticesByNotificationDate(ctx, date)
			if err != nil {
				s.logger.Error("get notification:" + err.Error())
			}

			sentNotifications := 0
			for _, notification := range notifications {
				notificationMarshaled, err := json.Marshal(notification)
				if err != nil {
					s.logger.Error("can't marshal notice:", "err", err.Error())
					continue
				}

				err = s.publisher.Publish(ctx, notificationMarshaled)
				if err != nil {
					s.logger.Error("can't publish notice:", "err", err.Error())
					continue
				}
				err = s.repository.UpdateIsNotified(ctx, notification.ID, 1)
				if err != nil {
					s.logger.Error("can't mark notice as sent:", "err", err.Error())
					continue
				}
				sentNotifications++
			}

			s.logger.Info("successfully publish:", "notifications", sentNotifications)

			s.truncateOlderEvents(ctx)
		}()

		select {
		case <-ctx.Done():
			s.logger.Info("Stop scheduler by context...")
			return
		case tm := <-ticker.C:
			s.logger.Info("check notifications on:", "time", tm)
		}
	}
}

func (s *Scheduler) truncateOlderEvents(ctx context.Context) {
	current := time.Now()
	yearAgo := current.AddDate(-1, 0, 0)

	err := s.repository.TruncateOlderEvents(ctx, yearAgo)
	if err != nil {
		s.logger.Error("can't truncate older events", "err", err.Error())
	}
}
