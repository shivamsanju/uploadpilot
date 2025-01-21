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
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

func Init() (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(LoggerMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Add companion proxy
	companionHandler, err := getCompanionProxyHandler(config.CompanionEndpoint, config.SelfEndpoint)
	if err != nil {
		return nil, err
	}
	router.Mount("/remote", companionHandler)

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Mount("/", Routes())
	})

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.WebServerPort),
	}

	return srv, nil
}

func getCompanionProxyHandler(companionEndpoint, selfEndpoint string) (http.Handler, error) {
	targetURL, err := url.Parse(companionEndpoint)
	if err != nil {
		infra.Log.Errorf("Failed to parse target URL: %v", err)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ModifyResponse = func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Location"), companionEndpoint) {
			resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), companionEndpoint, selfEndpoint, 1))
		}
		if strings.Contains(resp.Header.Get("Location"), url.QueryEscape(companionEndpoint)) {
			resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), url.QueryEscape(companionEndpoint), url.QueryEscape(selfEndpoint), 1))
		}
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
