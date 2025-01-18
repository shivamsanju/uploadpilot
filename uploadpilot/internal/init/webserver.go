package init

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
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/web"
)

func initWebServer(config *config.Config) error {
	g.RootPassword = config.RootPassword
	g.TusUploadDir = "./tmp"
	g.TusUploadBasePath = "/upload"
	g.FrontendURI = config.FrontendURI

	// Create the router and add the common middlewares.
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(web.LoggerMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Add upload pilot routes
	companionHandler, err := getCompanionProxyHandler(config.CompanionEndpoint, config.SelfEndpoint)
	if err != nil {
		return err
	}
	router.Mount("/remote", companionHandler)
	router.Group(func(r chi.Router) {
		r.Use(web.CorsMiddleware(config.FrontendURI))
		r.Mount("/", web.Routes())
	})

	// Start the web server.
	g.Log.Infof("starting webserver on port %d", config.WebServerPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.WebServerPort), router)
	if err != nil {
		g.Log.Errorf("failed to start webserver: %+v", err)
		return err
	}
	return nil
}

func getCompanionProxyHandler(companionEndpoint, selfEndpoint string) (http.Handler, error) {
	targetURL, err := url.Parse(companionEndpoint)
	if err != nil {
		g.Log.Errorf("Failed to parse target URL: %v", err)
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
