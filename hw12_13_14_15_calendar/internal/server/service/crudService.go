package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name IEventCrud --dir ./ --output ./../../../internal/mocks
type (
	IEventCrud interface {
		CreateEvent(ctx context.Context, event model.Event) (int64, error)
		UpdateEvent(ctx context.Context, event model.Event) (int64, error)
		DeleteEvent(ctx context.Context, id model.EventID) (int64, error)
		EventsByPeriodForOwner(ctx context.Context, ownerID model.OwnerID, start, end time.Time) ([]model.Event, error)
	}
)

type CalendarServer struct {
	eventpb.UnimplementedCalendarServer

	crud IEventCrud
}

func NewCalendarServer(crud IEventCrud) *CalendarServer {
	return &CalendarServer{
		crud: crud,
	}
}

func (cs *CalendarServer) CreateEvent(
	ctx context.Context,
	request *eventpb.CreateEventRequest,
) (*eventpb.CreateEventResponse, error) {
	uuid, err := cs.crud.CreateEvent(ctx, FromEvent(request.Event))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"%v : %s",
			storage.ErrDateBusy, request.Event.Start.AsTime(),
		)
	}
	if err != nil {
		return nil, err
	}

	return &eventpb.CreateEventResponse{InsertedUID: uuid}, nil
}

func (cs *CalendarServer) UpdateEvent(
	ctx context.Context,
	request *eventpb.UpdateEventRequest,
) (*eventpb.UpdateEventResponse, error) {
	updatedUID, err := cs.crud.UpdateEvent(ctx, FromEvent(request.Event))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, "date %s already busy", request.Event.Start.AsTime())
	}
	if err != nil {
		return nil, err
	}

	return &eventpb.UpdateEventResponse{
		Updated: updatedUID,
	}, nil
}

func (cs *CalendarServer) DeleteEvent(
	ctx context.Context,
	request *eventpb.DeleteEventRequest,
) (*eventpb.DeleteEventResponse, error) {
	deletedUID, err := cs.crud.DeleteEvent(ctx, model.EventID(request.Id))
	if err != nil {
		return nil, err
	}

	return &eventpb.DeleteEventResponse{
		Deleted: deletedUID,
	}, nil
}

func (cs *CalendarServer) EventsByPeriodAndOwner(
	ctx context.Context,
	request *eventpb.EventListRequest,
) (*eventpb.EventListResponse, error) {
	events, err := cs.crud.EventsByPeriodForOwner(
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
