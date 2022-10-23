package model

import "time"

type Event struct {
	ID               EventID
	Title            string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	OwnerID          OwnerID
	NotificationDate time.Time
	IsNotified       byte
}
