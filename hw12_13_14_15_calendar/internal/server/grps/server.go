package grps

import (
	"fmt"
	"net"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	address string
	logger  logger.ILogger
	config  *config.Config
}

func New(calendarServer eventpb.CalendarServer, logger logger.ILogger, config *config.Config) *Server {
	zapLogger := logger.GetInstance()
	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)
	// Create a server, make sure we put the grpc_ctxtags context before everything else.
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
		address: net.JoinHostPort(config.GRPS.Host, config.GRPS.Port),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
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
