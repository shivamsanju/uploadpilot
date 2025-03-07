package middlewares

import (
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func (m *Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// for uploader routes all origins are allowed
		if isCreateUploadRoute(r.URL.Path) || isFinishUploadRoute(r.URL.Path) {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Access-Control-Allow-Credentials", "true")
			response.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(
				append([]string{"Content-Type", "X-Api-Key", "X-Tenant-Id"}, supertokens.GetAllCORSHeaders()...),
				",",
			))
		} else {
			if slices.Contains(m.allowedOrigins, origin) {
				response.Header().Set("Access-Control-Allow-Origin", origin) // Dynamically set origin
				response.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(
				append([]string{"Content-Type", "X-Api-Key", "X-Tenant-Id"}, supertokens.GetAllCORSHeaders()...),
				",",
			))
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			response.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(response, r)
	})
}

func isCreateUploadRoute(path string) bool {
	pattern := `^/tenants/[^/]+/workspaces/[^/]+/uploads$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(path)
}

func isFinishUploadRoute(path string) bool {
	pattern := `^/tenants/[^/]+/workspaces/[^/]+/uploads/[^/]+/finish$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(path)
}
