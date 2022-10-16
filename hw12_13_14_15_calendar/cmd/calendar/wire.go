//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/grps"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/service"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/factory"
	"github.com/google/wire"
)

func InitializeDIForServer(config *config.Config, logger logger.ILogger) (*server.Server, error) {
	wire.Build(
		wire.Bind(new(service.IEventCrud), new(*repository.EventCrud)),
		wire.Bind(new(eventpb.CalendarServer), new(*service.CalendarServer)),
		service.NewCalendarServer,
		factory.MakeStorage,
		repository.NewEventCrud,
		rest.NewHandler,
		rest.New,
		grps.New,
		server.NewServerAggregator,
	)
	return &server.Server{}, nil
}
