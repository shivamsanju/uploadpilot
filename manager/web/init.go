package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/svc"
)

func InitWebServer(services *svc.Services) (*http.Server, error) {
	router := chi.NewRouter()
	appMiddlewares := NewAppMiddlewares(services.AuthService)

	// App middlewares
	router.Use(appMiddlewares.RecoveryMiddleware)
	router.Use(appMiddlewares.RequestIDMiddleware)
	router.Use(appMiddlewares.LoggerMiddleware)
	router.Use(appMiddlewares.CorsMiddleware)
	router.Use(appMiddlewares.RequestTimeoutMiddleware(30 * time.Second))

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Mount("/", Routes(services, appMiddlewares))
	})

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	return srv, nil
}
