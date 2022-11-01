package server

import (
	"context"
	"os"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
)

type Server struct {
	GRPC   *grpc.Server
	HTTP   *rest.Server
	logger logger.ILogger
}

func NewServerAggregator(grpcServer *grpc.Server, httpServer *rest.Server, logger logger.ILogger) *Server {
	return &Server{
		GRPC:   grpcServer,
		HTTP:   httpServer,
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context, stop context.CancelFunc) {
	s.logger.Info("calendar is running...")

	s.startHTTP(ctx, stop)

	s.startGRPC(ctx, stop)
}

func (s *Server) startGRPC(ctx context.Context, stop context.CancelFunc) {
	go func() {
		if err := s.GRPC.Start(ctx); err != nil {
			s.logger.Error("failed to start grpс server: " + err.Error())
			stop()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		s.GRPC.Stop()
	}()
}

func (s *Server) startHTTP(ctx context.Context, stop context.CancelFunc) {
	go func() {
		if err := s.HTTP.Start(ctx); err != nil {
			s.logger.Error("failed to start http server: " + err.Error())
			stop()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := s.HTTP.Stop(ctx); err != nil {
			s.logger.Error("failed to stop http server: " + err.Error())
		}
	}()
}
