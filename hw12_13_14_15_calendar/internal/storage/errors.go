package storage

import "errors"

var (
	ErrCantAddEvent         = errors.New("can't add event")
	ErrCantUpdateEvent      = errors.New("can't update event")
	ErrCantDeleteEvent      = errors.New("can't delete event")
	ErrCantConnectToStorage = errors.New("can't connect to storage")
	ErrNotFound             = errors.New("item not found")
	ErrDateBusy             = errors.New("date is already busy")
	ErrCantCreateStorage    = errors.New("can't create storage")
	ErrCantIdentifyStorage  = errors.New("can't identify storage")
)
