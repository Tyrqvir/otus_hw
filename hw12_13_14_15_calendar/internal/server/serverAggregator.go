package server

import (
	"context"
	"os"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	grpc "github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/grps"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
)

//nolint:lll
//go:generate protoc -I ../../api --go_out=../../api/eventpb --go_opt=paths=source_relative --go-grpc_out=../../api/eventpb ../../api/EventService.proto
//go:generate protoc -I ../../api --grpc-gateway_out=../../api/eventpb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true ../../api/EventService.proto
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

	s.startHttp(ctx, stop)

	s.startGRPS(ctx, stop)
}

func (s *Server) startGRPS(ctx context.Context, stop context.CancelFunc) {
	go func() {
		if err := s.GRPC.Start(); err != nil {
			s.logger.Error("failed to start grps server: " + err.Error())
			stop()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		s.GRPC.Stop()
	}()
}

func (s *Server) startHttp(ctx context.Context, stop context.CancelFunc) {
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
