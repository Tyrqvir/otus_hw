package memorystorage

import (
	"sync"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

type Storage struct {
	store map[string]model.Event
	mu    sync.RWMutex
}

func (s *Storage) AddEvent(event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, eventID)

	return nil
}

func (s *Storage) DailyEvents(date time.Time) ([]model.Event, error) {
	var events []model.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Start, 0)

		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() && eventDate.Day() == date.Day() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) WeeklyEvents(date time.Time) ([]model.Event, error) {
	var events []model.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Start, 0)
		eventYear, eventWeek := eventDate.ISOWeek()
		currentYear, currentWeek := date.ISOWeek()

		if eventYear == currentYear && eventWeek == currentWeek {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) MonthEvents(date time.Time) ([]model.Event, error) {
	var events []model.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Start, 0)
		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() {
			events = append(events, event)
		}
	}

	return events, nil
}

func New() *Storage {
	return &Storage{store: map[string]model.Event{}}
}
