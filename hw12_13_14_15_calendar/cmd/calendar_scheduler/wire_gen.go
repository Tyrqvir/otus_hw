// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/factory"
)

// Injectors from wire.go:

func InitializeDIForScheduler(config2 *config.Config, logger2 logger.Logger, rmqConfig broker.RMQConfig) (*scheduler.Scheduler, error) {
	connection := broker.NewConnection(rmqConfig, logger2, config2)
	producer := broker.NewProducer(config2, logger2, connection)
	eventRepository, err := factory.MakeStorage(config2)
	if err != nil {
		return nil, err
	}
	schedulerScheduler := scheduler.New(config2, logger2, producer, connection, eventRepository)
	return schedulerScheduler, nil
}
