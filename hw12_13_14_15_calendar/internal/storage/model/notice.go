package model

import "time"

type Notice struct {
	ID       EventID
	Title    string
	Datetime time.Time
	OwnerID  OwnerID
}
