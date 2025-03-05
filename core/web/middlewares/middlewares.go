package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/web/webutils"
)

type Middlewares struct {
	apiKeySvc      *services.APIKeyService
	allowedOrigins []string
}

func NewAppMiddlewares(apiKeySvc *services.APIKeyService, allowedOrigins []string) *Middlewares {
	return &Middlewares{
		apiKeySvc:      apiKeySvc,
		allowedOrigins: allowedOrigins,
	}
}

func (m *Middlewares) RecoveryMiddleware(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func (m *Middlewares) RequestIDMiddleware(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func (m *Middlewares) RequestTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

func (m *Middlewares) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			if ww.Status() >= 400 {
				log.Error().
					Str("request_id", middleware.GetReqID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("protocol", r.Proto).
					Str("remote_addr", r.RemoteAddr).
					Str("status", webutils.GetStatusLabel(ww.Status())).
					Int("bytes_written", int(ww.BytesWritten())).
					Int64("time_taken", time.Since(t1).Milliseconds()).
					Msg("request failed")
			} else {
				log.Info().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("bytes_written", int(ww.BytesWritten())).
					Int64("time_taken", time.Since(t1).Milliseconds()).
					Msg("request completed")
			}
		}()
		next.ServeHTTP(ww, r)
	})
}
