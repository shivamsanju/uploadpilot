package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/phuslu/log"
)

func AllowAllCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		headerKeys := r.Header.Get("Access-Control-Request-Headers")
		response.Header().Set("Access-Control-Allow-Origin", origin)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", headerKeys)

		// Handle preflight requests properly
		if r.Method == "OPTIONS" {
			response.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
			return
		}

		next.ServeHTTP(response, r)
	})
}

func AllowCorsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Origin", "X-Api-Key", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers

	}).Handler(next)
}

func LoggerMiddleware(next http.Handler) http.Handler {
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
					Str("status", GetStatusLabel(ww.Status())).
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
