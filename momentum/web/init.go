package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/momentum/internal/config"
)

func Init() (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(LoggerMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	h := Newhandler()

	router.Get("/", h.HealthCheck)
	router.Post("/trigger", h.TriggerWorkflow)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	infra.Log.Infof("Starting server on port %d", config.Port)

	return srv, nil
}
