package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

type IEventRepository interface {
	CreateEvent(ctx context.Context, e model.Event) (model.EventID, error)
	UpdateEvent(ctx context.Context, e model.Event) (int64, error)
	DeleteEvent(ctx context.Context, eID model.EventID) (int64, error)
	EventsByPeriodForOwner(ctx context.Context, ownerID model.OwnerID, start, end time.Time) ([]model.Event, error)
}

type EventCrud struct {
	eventRepository IEventRepository
}

func NewEventCrud(eventRepository IEventRepository) *EventCrud {
	return &EventCrud{
		eventRepository: eventRepository,
	}
}

func (ec *EventCrud) CreateEvent(ctx context.Context, event model.Event) (int64, error) {
	insertedID, err := ec.eventRepository.CreateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("cannot create event: %w", err)
	}

	return int64(insertedID), nil
}

func (ec *EventCrud) UpdateEvent(ctx context.Context, event model.Event) (int64, error) {
	updatedUID, err := ec.eventRepository.UpdateEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	return updatedUID, nil
}

func (ec *EventCrud) DeleteEvent(ctx context.Context, id model.EventID) (int64, error) {
	return ec.eventRepository.DeleteEvent(ctx, id)
}

func (ec *EventCrud) EventsByPeriodForOwner(
	ctx context.Context,
	ownerID model.OwnerID,
	start, end time.Time,
) ([]model.Event, error) {
	return ec.eventRepository.EventsByPeriodForOwner(ctx, ownerID, start, end)
}