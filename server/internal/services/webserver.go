package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shivamsanju/uploader/internal/config"
	"github.com/shivamsanju/uploader/internal/web"
	g "github.com/shivamsanju/uploader/pkg/globals"
)

func initWebServer(config *config.Config) error {
	g.TusUploadDir = "/tmp"

	// Create a new router with support for CORS and logging.
	router := chi.NewRouter()
	router.Use(web.CorsMiddleware(config.FrontendURI))
	router.Use(middleware.RequestID)
	router.Use(web.LoggerMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Mount the routes for the web app.
	r := web.Routes()
	router.Mount("/", r)

	g.RootPassword = config.RootPassword

	// Start the web server.
	g.Log.Infof("starting webserver on port %d", config.WebServerPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.WebServerPort), router)
	if err != nil {
		g.Log.Errorf("failed to start webserver: %+v", err)
	}

	return err
}
