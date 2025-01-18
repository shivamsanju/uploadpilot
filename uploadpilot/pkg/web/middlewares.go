package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/pkg/auth"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/utils"
)

func CorsMiddleware(frontendUri string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
			response.Header().Set("Access-Control-Allow-Origin", frontendUri)
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
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			g.Log.Infof(`id: %s | %s %s %s | %s | %s | %dB | %s`,
				middleware.GetReqID(r.Context()),  // RequestID (if set)
				r.Method,                          // Method
				r.URL.Path,                        // Path
				r.Proto,                           // Protocol
				r.RemoteAddr,                      // RemoteAddr
				utils.GetStatusLabel(ww.Status()), // "200 OK"
				ww.BytesWritten(),                 // Bytes Written
				time.Since(t1),                    // Elapsed
			)
		}()
		next.ServeHTTP(ww, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getBearerToken(r)
		if err != nil || len(token) == 0 {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}
		r.Header.Set("userId", claims.UserID)
		next.ServeHTTP(w, r)
	})
}

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}
