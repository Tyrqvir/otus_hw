package internalhttp

import (
	"encoding/json"
	"net/http"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/app"
)

func helloAction(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("hello-world")
	if err != nil {
		return
	}
}

func decorateHandler(logger app.Logger) http.HandlerFunc {
	handler := http.HandlerFunc(helloAction)
	handler = headersMiddleware(handler)
	handler = loggingMiddleware(handler, logger)

	return handler
}
