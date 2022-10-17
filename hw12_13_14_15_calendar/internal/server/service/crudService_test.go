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
		event := &eventpb.Event{}
		ctx := context.Background()
		insertedUID := int64(1)

		eventCRUD.On("CreateEvent", ctx, FromEvent(event)).Return(insertedUID, nil)

		server := NewCalendarServer(eventCRUD)
		response, err := server.CreateEvent(ctx, &eventpb.CreateEventRequest{Event: event})

		require.NoError(t, err)
		require.Equal(t, insertedUID, response.InsertedUid)
	})

	t.Run("error", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()
		event := &eventpb.Event{}

		eventCRUD.On("CreateEvent", ctx, FromEvent(event)).Return(int64(-1), fmt.Errorf("internal error"))

		server := NewCalendarServer(eventCRUD)
		response, err := server.CreateEvent(ctx, &eventpb.CreateEventRequest{Event: event})

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

		deletedUID := int64(1)

		eventCRUD.On("DeleteEvent", ctx, model.EventID(1)).Return(deletedUID, fmt.Errorf("error"))

		server := NewCalendarServer(eventCRUD)
		_, err := server.DeleteEvent(ctx, &eventpb.DeleteEventRequest{Id: int64(1)})

		require.Error(t, err)
	})
}

func TestCalendarServer_UpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventCRUD := &mocks.IEventCrud{}
		ctx := context.Background()
		event := &eventpb.Event{}
		updatedUID := int64(1)

		eventCRUD.On("UpdateEvent", ctx, FromEvent(event)).Return(updatedUID, nil)

		server := NewCalendarServer(eventCRUD)
		_, err := server.UpdateEvent(ctx, &eventpb.UpdateEventRequest{Event: event})

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
				ID:      1,
				OwnerID: 1,
				Start:   startData,
				End:     endData,
			},
			{
				ID:      2,
				OwnerID: 1,
				Start:   startData,
				End:     endData,
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
