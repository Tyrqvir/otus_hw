package app

import (
	"fmt"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Debug(string, ...interface{})
	Warn(string, ...interface{})
}

type Storage interface {
	AddEvent(e model.Event) error
	UpdateEvent(e model.Event) error
	DeleteEvent(eID string) error
	DailyEvents(date time.Time) ([]model.Event, error)
	WeeklyEvents(date time.Time) ([]model.Event, error)
	MonthEvents(date time.Time) ([]model.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) GetLogger() Logger {
	return a.logger
}

func (a *App) AddEvent(title string, description string, start int64, end int64, ownerID string, NotificationBefore int64) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("can't create uuid, %w", err)
	}

	event := model.Event{
		ID: id.String(), Title: title, Start: start, End: end, Description: description, OwnerID: ownerID, NotificationBefore: NotificationBefore,
	}

	return a.storage.AddEvent(event)
}

func (a *App) UpdateEvent(id, title string, description string, start int64, end int64, ownerID string, NotificationBefore int64) error {
	event := model.Event{
		ID: id, Title: title, Start: start, End: end, Description: description, OwnerID: ownerID, NotificationBefore: NotificationBefore,
	}

	return a.storage.UpdateEvent(event)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) DailyEvents(date time.Time) ([]model.Event, error) {
	return a.storage.DailyEvents(date)
}

func (a *App) WeeklyEvents(date time.Time) ([]model.Event, error) {
	return a.storage.WeeklyEvents(date)
}

func (a *App) MonthEvents(date time.Time) ([]model.Event, error) {
	return a.storage.MonthEvents(date)
}
