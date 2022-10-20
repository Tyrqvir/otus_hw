package model

type Event struct {
	ID                 string
	Title              string
	Start              int64
	End                int64
	Description        string
	OwnerID            string
	NotificationBefore int64
}
