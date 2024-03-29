package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pkg/errors"
)

type Server struct {
	httpServer http.Server
	logger     logger.Logger
}

func New(h http.Handler, logger logger.Logger, config *config.Config) *Server {
	return &Server{
		logger: logger,
		httpServer: http.Server{
			Addr:              net.JoinHostPort(config.HTTP.Host, config.HTTP.Port),
			Handler:           h,
			ReadTimeout:       config.HTTP.ReadTimeout,
			WriteTimeout:      config.HTTP.WriteTimeout,
			ReadHeaderTimeout: config.HTTP.ReadHeaderTimeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("start http server...")

	listenConfig := net.ListenConfig{}

	listener, err := listenConfig.Listen(ctx, "tcp", s.httpServer.Addr)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("can't start lister http server, %w", err)
	}

	err = s.httpServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("can't start serve http server, %w", err)
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stop http server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("can't shutdown http server, %w", err)
	}
	return nil
}
