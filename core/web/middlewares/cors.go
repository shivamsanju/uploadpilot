package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func (m *Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Validate and set allowed origins
		if origin != "" && slices.Contains(m.allowedOrigins, origin) {
			response.Header().Set("Access-Control-Allow-Origin", origin) // Dynamically set origin
			response.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", strings.Join(
			append([]string{"Content-Type", "X-Api-Key", "X-Tenant-Id"}, supertokens.GetAllCORSHeaders()...),
			",",
		))

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			response.WriteHeader(http.StatusNoContent) // No content needed for preflight
			return                                     // Return early to avoid calling next handler
		}

		next.ServeHTTP(response, r)
	})
}
