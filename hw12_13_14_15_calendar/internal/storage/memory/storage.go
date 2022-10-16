package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

type Storage struct {
	store         map[model.EventID]model.Event
	mu            sync.RWMutex
	incrementedID model.EventID
}

func (s *Storage) CreateEvent(_ context.Context, event model.Event) (model.EventID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, storeEvent := range s.store {
		if storeEvent.Start.Equal(event.Start) && storeEvent.OwnerID == event.OwnerID {
			return 0, storage.ErrDateBusy
		}
	}

	s.incrementedID++

	s.store[s.incrementedID] = event

	return s.incrementedID, nil
}

func (s *Storage) UpdateEvent(_ context.Context, event model.Event) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[event.ID] = event

	return 1, nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID model.EventID) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, eventID)

	return 1, nil
}

func (s *Storage) EventsByPeriodForOwner(
	_ context.Context,
	ownerID model.OwnerID,
	start, end time.Time,
) ([]model.Event, error) {
	var events []model.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		if ownerID == event.OwnerID && start.Before(event.Start) && end.After(event.End) {
			events = append(events, event)
		}
	}

	return events, nil
}

func New() *Storage {
	return &Storage{store: map[model.EventID]model.Event{}}
}
