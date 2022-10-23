package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/jmoiron/sqlx"
	//nolint:blank-imports
	_ "github.com/lib/pq"
)

const driver = "postgres"

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) (model.EventID, error) {
	result, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO events (
                    id,
                    title,
                    start_date,
                    end_date,
                    description,
                    owner_id,
                    notification_date
                    )
				VALUES (
				        :id,
				        :title,
				        :start_date,
				        :end_date,
				        :description,
				        :owner_id,
				        :notification_date
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

func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) (int64, error) {
	_, err := s.db.NamedExecContext(
		ctx,
		`UPDATE events
				SET
				    title=:title,
				    start_date=:startDate,
				    end_date=:endDate,
				    description=:description,
				    notification_date=:notification_date
              	WHERE
              	    id = :id`, &event)
	if err != nil {
		return 1, fmt.Errorf("%v, %w", storage.ErrCantUpdateEvent, err)
	}
	return 0, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID model.EventID) (int64, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1", eventID)

	if count, err := result.RowsAffected(); count == 0 {
		if err != nil {
			return 0, fmt.Errorf("can't get rows, %w", sql.ErrNoRows)
		}

		return 0, fmt.Errorf("events: %w", storage.ErrNotFound)
	}

	return 1, fmt.Errorf("%v, %w", storage.ErrCantDeleteEvent, err)
}

func (s *Storage) EventsByPeriodForOwner(
	ctx context.Context,
	ownerID model.OwnerID,
	startDate, endDate time.Time,
) ([]model.Event, error) {
	var result []model.Event

	err := s.db.SelectContext(
		ctx,
		&result,
		`SELECT
    				id, title, start_date, end_date, owner_id, description
				FROM
				    events
				WHERE
				    owner_id = $1 AND start_date BETWEEN $2 AND $3`,
		ownerID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) TruncateOlderEvents(ctx context.Context, date time.Time) error {
	query := `DELETE FROM events WHERE start_date <= $1 AND is_notified = $2;`
	result, err := s.db.ExecContext(ctx, query, date, 1)
	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 {
		if err != nil {
			return fmt.Errorf("can't get rows, %w", sql.ErrNoRows)
		}
	}

	return nil
}

func (s *Storage) NoticesByNotificationDate(ctx context.Context, date time.Time) ([]model.Notice, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, title, notification_date, owner_id FROM events WHERE is_notified = $1 AND notification_date <= $2`,
		0,
		date,
	)
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []model.Notice

	for rows.Next() {
		var notice model.Notice
		if err := rows.Scan(
			&notice.ID,
			&notice.Title,
			&notice.Datetime,
			&notice.OwnerID,
		); err != nil {
			return nil, err
		}

		notices = append(notices, notice)
	}

	return notices, nil
}

func (s *Storage) UpdateIsNotified(ctx context.Context, id model.EventID, isNotified byte) error {
	query := `UPDATE events SET is_notified = $1 WHERE id = $2;`
	_, err := s.db.ExecContext(ctx, query, isNotified, id)
	if err != nil {
		return fmt.Errorf("%v, %w", storage.ErrCantUpdateEvent, err)
	}

	return nil
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
