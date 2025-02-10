package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/momentum/internal/config"
	"github.com/uploadpilot/uploadpilot/momentum/internal/utils"
)

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
		if r.Header.Get("Secret") != config.SecretKey {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
