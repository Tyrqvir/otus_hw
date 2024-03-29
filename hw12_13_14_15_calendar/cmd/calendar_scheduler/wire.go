//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/factory"
	"github.com/google/wire"
)

func InitializeDIForScheduler(
	config *config.Config,
	logger logger.Logger,
	rmqConfig broker.RMQConfig,
) (*scheduler.Scheduler, error) {
	wire.Build(
		factory.MakeStorage,
		broker.NewConnection,
		broker.NewProducer,
		scheduler.New,
	)

	return &scheduler.Scheduler{}, nil
}
