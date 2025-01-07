package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"github.com/shivamsanju/uploader/web/utils"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func CorsMiddleware(frontendUri string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
			response.Header().Set("Access-Control-Allow-Origin", frontendUri)
			response.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
				response.Header().Set("Access-Control-Allow-Methods", "*")
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
		session.VerifySession(nil, next.ServeHTTP).ServeHTTP(w, r)
	})
}
