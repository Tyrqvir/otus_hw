package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger/zap"
	internalhttp "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	storageFactory "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/storage/factory"
)

const versionArgKey = "version"

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == versionArgKey {
		printVersion()
		return
	}

	fmt.Println(configFile)

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	logg := zap.New(cfg.Logger.Level)

	storage, err := storageFactory.MakeStorage(cfg)
	if err != nil {
		panic(err)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
