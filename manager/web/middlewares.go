package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/phuslu/log"
	"github.com/uploadpilot/manager/internal/auth"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/utils"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", config.FrontendURI)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Header().Set("Access-Control-Allow-Headers", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
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
					Str("status", utils.GetStatusLabel(ww.Status())).
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify api key
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey == config.APIKey {
			ctx := context.WithValue(r.Context(), dto.UserIDContextKey, "api-key")
			ctx = context.WithValue(ctx, dto.EmailContextKey, "api-key")
			ctx = context.WithValue(ctx, dto.NameContextKey, "api-key")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// verify jwt
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		token := parts[1]
		if len(token) == 0 {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, dto.EmailContextKey, claims.Email)
		ctx = context.WithValue(ctx, dto.NameContextKey, claims.Name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
