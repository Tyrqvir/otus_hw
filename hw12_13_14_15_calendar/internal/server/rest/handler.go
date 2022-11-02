package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/api/eventpb"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewHandler(config *config.Config, logger logger.ILogger) (http.Handler, error) {
	gw := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := eventpb.RegisterCalendarHandlerFromEndpoint(
		context.Background(),
		gw,
		config.GRPC.Port,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("register calendar service failed: %w", err)
	}

	gwmux := http.NewServeMux()
	handler := loggingMiddleware(gw, logger)
	gwmux.Handle("/", handler)

	return gwmux, nil
}
