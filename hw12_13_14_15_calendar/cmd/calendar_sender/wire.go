//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/sender"
	"github.com/google/wire"
)

func InitializeDIForSender(
	config *config.Config,
	logger logger.Logger,
	rmqConfig broker.RMQConfig,
) (*sender.Sender, error) {
	wire.Build(
		broker.NewConnection,
		broker.NewConsumer,
		sender.New,
	)

	return &sender.Sender{}, nil
}
