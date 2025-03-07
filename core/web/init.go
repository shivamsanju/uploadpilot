package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/core/config"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/web/middlewares"
	"github.com/uploadpilot/core/web/routes"
)

func NewWebserver(appConfig *config.Config, services *services.Services) (*http.Server, error) {
	router := chi.NewRouter()
	appMiddlewares := middlewares.NewAppMiddlewares(services.WorkspaceService, services.APIKeyService, appConfig.AllowedOrigins)

	// App middlewares
	router.Use(appMiddlewares.RecoveryMiddleware)
	router.Use(appMiddlewares.RequestIDMiddleware)
	router.Use(appMiddlewares.LoggerMiddleware)
	router.Use(appMiddlewares.CorsMiddleware)
	router.Use(appMiddlewares.RequestTimeoutMiddleware(30 * time.Second))

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Mount("/", routes.NewAppRoutesV1(services, appMiddlewares))
	})

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", appConfig.Port),
	}

	return srv, nil
}
