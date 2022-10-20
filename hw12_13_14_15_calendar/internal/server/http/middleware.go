package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tyrqvir/otus_hw/hw12_13_14_15_calendar/internal/app"
)

func headersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

type responseWriterDecorator struct {
	http.ResponseWriter

	status int
}

func (rd *responseWriterDecorator) WriteHeader(status int) {
	rd.status = status
	rd.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.HandlerFunc, logger app.Logger) http.HandlerFunc {
	start := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		wd := &responseWriterDecorator{ResponseWriter: w}

		next.ServeHTTP(wd, r)

		logger.Debug(fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			start.Format(time.RFC1123Z),
			r.Method,
			r.URL.Path,
			r.Proto,
			wd.status,
			time.Since(start).Milliseconds(),
			r.UserAgent()))
	}
}
