package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/svc"
)

func InitWebServer(services *svc.Services) (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(LoggerMiddleware)

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Mount("/", Routes(services))
	})

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	return srv, nil
}
