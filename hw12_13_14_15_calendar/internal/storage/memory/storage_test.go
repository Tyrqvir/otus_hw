package memorystorage

import (
	"testing"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/stretchr/testify/require"
)

const dummyID = "dummy id"

func dummyEvent() model.Event {
	return model.Event{
		ID:                 dummyID,
		Title:              "dummy title",
		Start:              time.Now().Unix(),
		End:                time.Now().AddDate(0, 1, 0).Unix(),
		Description:        "dummy description",
		OwnerID:            "test_owner_id",
		NotificationBefore: 0,
	}
}

func dummyStorage() *Storage {
	event := dummyEvent()
	dummyStorage := New()
	err := dummyStorage.AddEvent(event)
	if err != nil {
		panic("Can't create dummy storage")
	}

	return dummyStorage
}

func TestStorage(t *testing.T) {
	t.Run("test add event", func(t *testing.T) {
		storage := New()
		event := model.Event{}
		err := storage.AddEvent(event)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Len(t, storage.store, 1)
	})

	t.Run("test update event", func(t *testing.T) {
		storage := dummyStorage()
		require.Len(t, storage.store, 1)

		event := storage.store[dummyID]

		event.Title = "Updated title"
		event.Description = "Updated description"

		err := storage.UpdateEvent(event)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Len(t, storage.store, 1)

		require.Equal(t, "Updated description", storage.store[dummyID].Description)
		require.Equal(t, "Updated title", storage.store[dummyID].Title)
	})

	t.Run("test delete event", func(t *testing.T) {
		storage := dummyStorage()
		require.Len(t, storage.store, 1)

		event := storage.store[dummyID]

		err := storage.DeleteEvent(event.ID)
		require.NoError(t, err)
		require.Nil(t, err)

		require.Len(t, storage.store, 0)
	})

	t.Run("test lists", func(t *testing.T) {
		storage := New()
		require.Len(t, storage.store, 0)

		err := storage.AddEvent(dummyEvent())
		require.NoError(t, err)

		now := time.Now()

		daily, err := storage.DailyEvents(now)
		require.NoError(t, err)
		require.Len(t, daily, 1)

		weekly, err := storage.DailyEvents(now)
		require.NoError(t, err)
		require.Len(t, weekly, 1)

		monthly, err := storage.DailyEvents(now)
		require.NoError(t, err)
		require.Len(t, monthly, 1)
	})

	t.Run("test add 2 same events", func(t *testing.T) {
		storage := New()
		event1 := model.Event{ID: "event1"}
		event2 := model.Event{ID: "event2"}
		err := storage.AddEvent(event1)
		require.NoError(t, err)

		err = storage.AddEvent(event2)
		require.NoError(t, err)

		require.Len(t, storage.store, 2)
	})
}
