package model

import "time"

type Event struct {
	ID               EventID   `db:"id"`
	Title            string    `db:"title"`
	StartDate        time.Time `db:"start_date"`
	EndDate          time.Time `db:"end_date"`
	Description      string    `db:"description"`
	OwnerID          OwnerID   `db:"owner_id"`
	NotificationDate time.Time `db:"notification_date"`
	IsNotified       int64     `db:"is_notified"`
}
