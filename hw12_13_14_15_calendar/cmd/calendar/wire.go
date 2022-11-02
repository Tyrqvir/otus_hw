//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/service"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/factory"
	"github.com/google/wire"
)

func setupWire(config *config.Config, logger logger.Logger) (*server.Server, error) {
	wire.Build(
		wire.Bind(new(eventpb.CalendarServer), new(*service.CalendarServer)),
		service.NewCalendarServer,
		factory.MakeStorage,
		rest.NewHandler,
		rest.New,
		grpc.New,
		server.NewServerAggregator,
	)
	return &server.Server{}, nil
}
