package sqlstorage

import (
	"context"
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

func (s *Storage) CreateEvent(_ context.Context, event model.Event) (model.EventID, error) {
	result, err := s.db.NamedExec(
		`INSERT INTO events (
                    id,
                    title,
                    start,
                    end,
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
				        )`, &event)
	if err != nil {
		return 0, fmt.Errorf("%v, %w", storage.ErrCantAddEvent, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%v, %w", storage.ErrNotFound, err)
	}

	return model.EventID(id), nil
}

func (s *Storage) UpdateEvent(_ context.Context, event model.Event) (int64, error) {
	_, err := s.db.NamedExec(
		`UPDATE events
				SET
				    title=:title,
				    start=:start,
				    end=:end,
				    description=:description,
				    notification_before=:notification_before
              	WHERE
              	    id = :id`, &event)

	return 1, fmt.Errorf("%v, %w", storage.ErrCantUpdateEvent, err)
}

func (s *Storage) DeleteEvent(_ context.Context, eventID model.EventID) (int64, error) {
	result, err := s.db.Exec("DELETE FROM events WHERE id=$1", eventID)

	if count, err := result.RowsAffected(); count == 0 {
		if err != nil {
			return 0, fmt.Errorf("can't get rows, %w", sql.ErrNoRows)
		}

		return 0, fmt.Errorf("events: %w", storage.ErrNotFound)
	}

	return 1, fmt.Errorf("%v, %w", storage.ErrCantDeleteEvent, err)
}

func (s *Storage) EventsByPeriodForOwner(
	_ context.Context,
	ownerID model.OwnerID,
	start, end time.Time,
) ([]model.Event, error) {
	var result []model.Event

	err := s.db.Select(
		&result,
		`SELECT id, title, start, end, owner_id, description
				FROM events
				WHERE owner_id = $1 AND start BETWEEN $2 AND $3`,
		ownerID, start, end)
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
