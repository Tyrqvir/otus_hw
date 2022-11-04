package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

//nolint:lll
//go:generate protoc -I ../../../api --go_out=../../../api/eventpb --go_opt=paths=source_relative --go-grpc_out=../../../api/eventpb ../../../api/EventService.proto
//go:generate protoc -I ../../../api --grpc-gateway_out=../../../api/eventpb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true ../../../api/EventService.proto
type Server struct {
	server  *grpc.Server
	address string
	logger  logger.Logger
	config  *config.Config
}

func New(calendarServer eventpb.CalendarServer, logger logger.Logger, config *config.Config) *Server {
	zapLogger := logger.GetInstance()
	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger),
		),
	)

	eventpb.RegisterCalendarServer(server, calendarServer)

	return &Server{
		logger:  logger,
		server:  server,
		config:  config,
		address: config.GRPC.Port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listenConfig := net.ListenConfig{}

	listener, err := listenConfig.Listen(ctx, "tcp", s.address)
	if err != nil {
		return fmt.Errorf("start grpc server failed: %w", err)
	}

	s.logger.Info("start grpc server...")

	return s.server.Serve(listener)
}

func (s *Server) Stop() {
	s.logger.Info("stop grpc server...")

	s.server.GracefulStop()
}
