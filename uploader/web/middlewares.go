package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/uploader/internal/infra"
)

func AllowAllCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
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
				middleware.GetReqID(r.Context()), // RequestID (if set)
				r.Method,                         // Method
				r.URL.Path,                       // Path
				r.Proto,                          // Protocol
				r.RemoteAddr,                     // RemoteAddr
				GetStatusLabel(ww.Status()),      // "200 OK"
				ww.BytesWritten(),                // Bytes Written
				time.Since(t1),                   // Elapsed
			)
		}()
		next.ServeHTTP(ww, r)
	})
}
