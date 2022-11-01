package repository

import (
	"context"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

//go:generate mockery --name IEventRepository --dir ./ --output ./../../internal/mocks
type IEventRepository interface {
	CreateEvent(ctx context.Context, e model.Event) (model.EventID, error)
	UpdateEvent(ctx context.Context, e model.Event) (bool, error)
	DeleteEvent(ctx context.Context, eID model.EventID) (bool, error)
	EventsByPeriodForOwner(ctx context.Context, ownerID model.OwnerID, startDate, endDate time.Time) ([]model.Event, error)
	TruncateOlderEvents(ctx context.Context, date time.Time) error
	NoticesByNotificationDate(ctx context.Context, date time.Time) ([]model.Notice, error)
	UpdateIsNotified(ctx context.Context, id model.EventID, isNotified int64) error
}

type EventRepository struct {
	IEventRepository
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}

}
