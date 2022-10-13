package sqlstorage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/jmoiron/sqlx"
)

const driver = "postgres"

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) AddEvent(e model.Event) error {
	_, err := s.db.NamedExec(
		`INSERT INTO events (
                    id,
                    title,
                    start,
                    finish,
                    description,
                    owner_id,
                    notification_before
                    )
				VALUES (
				        :id,
				        :title,
				        :start,
				        :end,
				        :description,
				        :owner_id,
				        :notification_before
				        )`, &e)

	return fmt.Errorf("%v, %w", storage.ErrCantAddEvent, err)
}

func (s *Storage) UpdateEvent(e model.Event) error {
	_, err := s.db.NamedExec(
		`UPDATE events
				SET
				    title=:title,
				    start=:date,
				    finish=:end,
				    description=:description,
				    notification_before=:notification_before
              	WHERE
              	    id = :id`, &e)

	return fmt.Errorf("%v, %w", storage.ErrCantUpdateEvent, err)
}

func (s *Storage) DeleteEvent(eID string) error {
	res, err := s.db.Exec("DELETE FROM events WHERE id=$1", eID)

	if count, err := res.RowsAffected(); count == 0 {
		if err != nil {
			return fmt.Errorf("can't get rows, %w", sql.ErrNoRows)
		}

		return fmt.Errorf("events: %w", storage.ErrNotFound)
	}

	return fmt.Errorf("%v, %w", storage.ErrCantDeleteEvent, err)
}

func (s *Storage) DailyEvents(date time.Time) ([]model.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	from := start.Unix()
	to := start.AddDate(0, 0, 1).Unix()

	result, err := s.eventsByDate(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) WeeklyEvents(date time.Time) ([]model.Event, error) {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	week := date.Add(time.Duration(offset*24) * time.Hour)
	start := time.Date(week.Year(), week.Month(), week.Day(), 0, 0, 0, 0, week.Location())
	from := start.Unix()
	to := start.AddDate(0, 0, 7).Unix()

	result, err := s.eventsByDate(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) MonthEvents(date time.Time) ([]model.Event, error) {
	current := time.Date(date.Year(), date.Month(), 0, 0, 0, 0, 0, date.Location())
	from := current.Unix()
	to := current.AddDate(0, 1, 0).Unix()

	result, err := s.eventsByDate(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) eventsByDate(from int64, to int64) ([]model.Event, error) {
	var result []model.Event

	err := s.db.Select(&result, "SELECT * FROM events WHERE start BETWEEN $1 AND $2", from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func New(dataSourceName string) (*Storage, error) {
	db, err := sqlx.Connect(driver, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%v, %w", storage.ErrCantConnectToStorage, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
