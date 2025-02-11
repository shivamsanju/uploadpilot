package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
)

func Init() (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(LoggerMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Mount the uploadpilot web routes
	router.Group(func(r chi.Router) {
		r.Mount("/", Routes())
	})

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	return srv, nil
}
