package model

import "time"

type (
	EventID int64
	OwnerID int64
)

type Event struct {
	ID                 EventID
	Title              string
	Start              time.Time
	End                time.Time
	Description        string
	OwnerID            OwnerID
	NotificationBefore time.Time
}
