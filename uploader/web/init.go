package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploader/internal/service"
)

func InitWebserver(svc *service.Service) (*http.Server, error) {
	appConfig := config.GetAppConfig()
	h := Newhandler(svc)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.RequestID)
	router.Use(LoggerMiddleware)

	// Add companion proxy
	rp, err := companionReverseProxy(appConfig.CompanionEndpoint, appConfig.UploaderEndpoint)
	if err != nil {
		return nil, err
	}
	router.Mount("/remote", rp)

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Use(AllowAllCorsMiddleware)
		r.Get("/health", h.HealthCheck)
		r.Get("/config/{workspaceId}", h.GetUploaderConfig)
		r.Mount("/upload/{workspaceId}", h.GetUploadHandler())
	})

	// infra.Log.Infof("Starting server on port %d", appConfig.Port)
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", appConfig.Port),
	}

	return srv, nil
}

func companionReverseProxy(companionEndpoint, uploaderEndpoint string) (http.Handler, error) {
	targetURL, err := url.Parse(companionEndpoint)
	if err != nil {
		// infra.Log.Errorf("Failed to parse target URL: %v", err)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ModifyResponse = func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Location"), companionEndpoint) {
			resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), companionEndpoint, uploaderEndpoint, 1))
		}
		if strings.Contains(resp.Header.Get("Location"), url.QueryEscape(companionEndpoint)) {
			resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), url.QueryEscape(companionEndpoint), url.QueryEscape(uploaderEndpoint), 1))
		}
		// infra.Log.Infof("modified response: %s", resp.Header.Get("Location"))
		return nil
	}

	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = targetURL.Host
		proxy.ServeHTTP(w, r)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		proxyHandler.ServeHTTP(w, r)
	}), nil
}
