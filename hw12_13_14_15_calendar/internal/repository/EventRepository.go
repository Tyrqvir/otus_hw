package repository

import (
	"context"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

//go:generate mockery --name IEventRepository --dir ./ --output ./../../internal/mocks
type IEventRepository interface {
	CreateEvent(ctx context.Context, event model.Event) (model.EventID, error)
	UpdateEvent(ctx context.Context, event model.Event) (bool, error)
	DeleteEvent(ctx context.Context, id model.EventID) (bool, error)
	EventsByPeriodForOwner(ctx context.Context, ownerID model.OwnerID, start, end time.Time) ([]model.Event, error)
}

type EventRepository struct {
	IEventRepository
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}
