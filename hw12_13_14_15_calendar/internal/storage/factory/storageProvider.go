package factory

import (
	"fmt"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	SQLProvider    = "sql"
	MemoryProvider = "in-memory"
)

func MakeStorage(config *config.Config) (repository.EventRepository, error) {
	switch config.DB.Provider {
	case MemoryProvider:
		return memorystorage.New(), nil
	case SQLProvider:
		instance, err := sqlstorage.New(config.DB.DSN)
		if err != nil {
			return nil, fmt.Errorf("%v, %w", storage.ErrCantCreateStorage, err)
		}
		return instance, nil
	}
	return nil, storage.ErrCantIdentifyStorage
}
