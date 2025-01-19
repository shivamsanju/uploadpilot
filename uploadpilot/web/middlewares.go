package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/internal/auth"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
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
			infra.Log.Infof(`id: %s | %s %s %s | %s | %s | %dB | %s`,
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
		r.Header.Set("userId", claims.UserID)
		r.Header.Set("email", claims.Email)
		r.Header.Set("name", claims.Name)
		next.ServeHTTP(w, r)
	})
}
