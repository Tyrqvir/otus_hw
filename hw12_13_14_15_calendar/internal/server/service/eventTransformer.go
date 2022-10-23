package service

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEvent(event model.Event) *eventpb.Event {
	return &eventpb.Event{
		Id:               int64(event.ID),
		Title:            event.Title,
		Description:      event.Description,
		StartDate:        timestamppb.New(event.StartDate),
		EndDate:          timestamppb.New(event.EndDate),
		OwnerId:          int64(event.OwnerID),
		NotificationDate: timestamppb.New(event.NotificationDate),
	}
}

func FromEvent(event *eventpb.Event) model.Event {
	return model.Event{
		ID:               model.EventID(event.Id),
		Title:            event.Title,
		Description:      event.Description,
		StartDate:        event.StartDate.AsTime(),
		EndDate:          event.EndDate.AsTime(),
		OwnerID:          model.OwnerID(event.OwnerId),
		NotificationDate: event.NotificationDate.AsTime(),
	}
}

func ToEvents(events []model.Event) []*eventpb.Event {
	pbEvents := make([]*eventpb.Event, 0, len(events))

	for _, event := range events {
		pbEvents = append(pbEvents, ToEvent(event))
	}

	return pbEvents
}
