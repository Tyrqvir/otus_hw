package service

import (
	"context"
	"errors"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalendarServer struct {
	eventpb.UnimplementedCalendarServer

	repository repository.EventRepository
}

func NewCalendarServer(repository repository.EventRepository) *CalendarServer {
	return &CalendarServer{
		repository: repository,
	}
}

func (cs *CalendarServer) CreateEvent(
	ctx context.Context,
	request *eventpb.CreateEventRequest,
) (*eventpb.CreateEventResponse, error) {
	id, err := cs.repository.CreateEvent(ctx, FromEvent(request.CommonEvent))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"%v : %s",
			storage.ErrDateBusy, request.CommonEvent.StartDate.AsTime(),
		)
	}
	if err != nil {
		return nil, err
	}

	return &eventpb.CreateEventResponse{InsertedId: int64(id)}, nil
}

func (cs *CalendarServer) UpdateEvent(
	ctx context.Context,
	request *eventpb.UpdateEventRequest,
) (*eventpb.UpdateEventResponse, error) {
	_, err := cs.repository.UpdateEvent(ctx, FromEvent(request.CommonEvent))
	if errors.Is(err, storage.ErrDateBusy) {
		return &eventpb.UpdateEventResponse{
			Updated: false,
		}, status.Errorf(codes.InvalidArgument, "date %s already busy", request.CommonEvent.StartDate.AsTime())
	}
	if err != nil {
		return &eventpb.UpdateEventResponse{
			Updated: false,
		}, err
	}

	return &eventpb.UpdateEventResponse{
		Updated: true,
	}, nil
}

func (cs *CalendarServer) DeleteEvent(
	ctx context.Context,
	request *eventpb.DeleteEventRequest,
) (*eventpb.DeleteEventResponse, error) {
	_, err := cs.repository.DeleteEvent(ctx, model.EventID(request.Id))
	if err != nil {
		return nil, err
	}

	return &eventpb.DeleteEventResponse{
		Deleted: true,
	}, nil
}

func (cs *CalendarServer) EventsByPeriodAndOwner(
	ctx context.Context,
	request *eventpb.EventListRequest,
) (*eventpb.EventListResponse, error) {
	events, err := cs.repository.EventsByPeriodForOwner(
		ctx,
		model.OwnerID(request.Owner),
		request.Start.AsTime(),
		request.End.AsTime(),
	)
	if err != nil {
		return nil, err
	}

	return &eventpb.EventListResponse{Events: ToEvents(events)}, nil
}
