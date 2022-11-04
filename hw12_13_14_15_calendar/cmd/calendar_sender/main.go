package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	globalLogger "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/app/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logger := globalLogger.New(cfg.Logger.Level)

	brokerConf := broker.RMQConfig{
		Tag:            cfg.Consumer.Tag,
		ConnectionName: cfg.Consumer.ConnectionName,
	}

	sender, err := InitializeDIForSender(cfg, logger, brokerConf)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	go func() {
		sender.Handle(ctx)
	}()

	<-ctx.Done()

	logger.Info("shutting down")

	if err := sender.Shutdown(); err != nil {
		logger.Errorf("error during shutdown: %s", err)
	}
}
