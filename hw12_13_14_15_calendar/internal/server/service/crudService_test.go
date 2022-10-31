package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/mocks"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCalendarServer_CreateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		event := &eventpb.CommonEvent{}
		ctx := context.Background()
		insertedID := int64(1)

		eventCRUD.On("CreateEvent", ctx, FromEvent(event)).Return(insertedID, nil)

		server := NewCalendarServer(eventCRUD)
		response, err := server.CreateEvent(ctx, &eventpb.CreateEventRequest{CommonEvent: event})

		require.NoError(t, err)
		require.Equal(t, insertedID, response.InsertedId)
	})

	t.Run("error", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()
		event := &eventpb.CommonEvent{}

		eventCRUD.On("CreateEvent", ctx, FromEvent(event)).Return(int64(-1), fmt.Errorf("internal error"))

		server := NewCalendarServer(eventCRUD)
		response, err := server.CreateEvent(ctx, &eventpb.CreateEventRequest{CommonEvent: event})

		require.Error(t, err)
		require.Nil(t, response)
	})
}

func TestCalendarServer_DeleteEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()
		rows := int64(1)
		eventID := int64(1)

		eventCRUD.On("DeleteEvent", ctx, model.EventID(1)).Return(rows, nil)

		server := NewCalendarServer(eventCRUD)
		_, err := server.DeleteEvent(ctx, &eventpb.DeleteEventRequest{Id: eventID})

		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()

		deletedID := int64(1)

		eventCRUD.On("DeleteEvent", ctx, model.EventID(1)).Return(deletedID, fmt.Errorf("error"))

		server := NewCalendarServer(eventCRUD)
		_, err := server.DeleteEvent(ctx, &eventpb.DeleteEventRequest{Id: int64(1)})

		require.Error(t, err)
	})
}

func TestCalendarServer_UpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()
		event := &eventpb.CommonEvent{}
		updatedID := int64(1)

		eventCRUD.On("UpdateEvent", ctx, FromEvent(event)).Return(updatedID, nil)

		server := NewCalendarServer(eventCRUD)
		_, err := server.UpdateEvent(ctx, &eventpb.UpdateEventRequest{CommonEvent: event})

		require.NoError(t, err)
	})
}

func TestCalendarServer_EventsByPeriodAndOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}

		currentTime := timestamppb.Now()

		startData := currentTime.AsTime()
		endData := currentTime.AsTime().Add(5 * time.Minute)

		items := []model.Event{
			{
				ID:        1,
				OwnerID:   1,
				StartDate: startData,
				EndDate:   endData,
			},
			{
				ID:        2,
				OwnerID:   1,
				StartDate: startData,
				EndDate:   endData,
			},
		}

		ctx := context.Background()
		eventCRUD.On("EventsByPeriodForOwner", ctx, model.OwnerID(1), startData, endData).Return(items, nil)

		server := NewCalendarServer(eventCRUD)
		foundedEvents, err := server.EventsByPeriodAndOwner(ctx, &eventpb.EventListRequest{
			Owner: 1,
			Start: currentTime,
			End:   timestamppb.New(currentTime.AsTime().Add(5 * time.Minute)),
		})

		require.NoError(t, err)
		require.Len(t, foundedEvents.Events, 2)
	})
}
