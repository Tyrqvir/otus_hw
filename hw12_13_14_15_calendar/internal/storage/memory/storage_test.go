package memorystorage

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/stretchr/testify/require"
)

const (
	dummyID = model.EventID(1)
	ownerID = model.OwnerID(1)
)

var ctx = context.Background()

func dummyEvent() model.Event {
	return model.Event{
		ID:                 dummyID,
		Title:              "dummy title",
		Start:              time.Now(),
		End:                time.Now(),
		Description:        "dummy description",
		OwnerID:            ownerID,
		NotificationBefore: time.Now(),
	}
}

func dummyStorage() *Storage {
	event := dummyEvent()
	dummyStorage := New()
	_, err := dummyStorage.CreateEvent(ctx, event)
	if err != nil {
		log.Fatalln("Can't create dummy storage:", err)
	}

	return dummyStorage
}

func TestStorage(t *testing.T) {
	t.Run("test add event", func(t *testing.T) {
		storage := New()
		event := model.Event{}
		id, err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Len(t, storage.store, 1)
		require.Equal(t, id, dummyID)
	})

	t.Run("test update event", func(t *testing.T) {
		storage := dummyStorage()
		require.Len(t, storage.store, 1)

		event := storage.store[dummyID]

		event.Title = "Updated title"
		event.Description = "Updated description"

		_, err := storage.UpdateEvent(ctx, event)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Len(t, storage.store, 1)
		require.Equal(t, storage.store[dummyID].ID, dummyID)

		require.Equal(t, "Updated description", storage.store[dummyID].Description)
		require.Equal(t, "Updated title", storage.store[dummyID].Title)
	})

	t.Run("test delete event", func(t *testing.T) {
		storage := dummyStorage()
		require.Len(t, storage.store, 1)

		event := storage.store[dummyID]

		_, err := storage.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)
		require.Nil(t, err)

		require.Len(t, storage.store, 0)
	})

	t.Run("test lists", func(t *testing.T) {
		storage := New()
		now := time.Now()
		event1 := model.Event{ID: 1, Start: now.Add(-3 * time.Minute), End: now.Add(4 * time.Minute), OwnerID: ownerID}
		event2 := model.Event{ID: 2, Start: now.Add(-2 * time.Minute), End: now.Add(10 * time.Minute), OwnerID: ownerID}
		_, _ = storage.CreateEvent(ctx, event1)
		_, _ = storage.CreateEvent(ctx, event2)
		require.Len(t, storage.store, 2)

		list, err := storage.EventsByPeriodForOwner(ctx, ownerID, now.Add(-5*time.Minute), now.Add(5*time.Minute))
		require.NoError(t, err)
		require.Len(t, list, 1)
	})

	t.Run("test add 2 with same ID but different start date", func(t *testing.T) {
		storage := New()
		currentTime := time.Now()
		event1 := model.Event{ID: 1, Start: currentTime}
		event2 := model.Event{ID: 1, Start: currentTime.Add(5 * time.Minute)}
		_, err := storage.CreateEvent(ctx, event1)
		require.NoError(t, err)

		_, err = storage.CreateEvent(ctx, event2)
		require.NoError(t, err)

		require.Len(t, storage.store, 2)
	})

	t.Run("test add 2 with same Date and different ID", func(t *testing.T) {
		storage := New()
		currentTime := time.Now()
		event1 := model.Event{ID: 1, Start: currentTime}
		event2 := model.Event{ID: 2, Start: currentTime}
		_, err := storage.CreateEvent(ctx, event1)
		require.NoError(t, err)

		_, err = storage.CreateEvent(ctx, event2)
		require.Error(t, err)
	})
}
