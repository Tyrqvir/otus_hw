package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	globalLogger "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

const versionArgKey = "version"

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/app/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == versionArgKey {
		printVersion()
		return
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logger := globalLogger.New(cfg.Logger.Level)

	server, err := setupWire(cfg, logger)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	server.Run(ctx, stop)
	<-ctx.Done()
}
