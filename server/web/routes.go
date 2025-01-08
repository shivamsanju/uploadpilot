package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/shivamsanju/uploader/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/workflows", GetWorkflowRoutes())
	r.Mount("/storage/connectors", GetStorageConnectorRoutes())
	return r
}

func GetWorkflowRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewWorkflowHandler()
	r.Post("/", h.CreateWorkflow)
	r.Get("/", h.ListWorkflows)
	r.Get("/{id}", h.GetWorkflow)
	return r
}

func GetStorageConnectorRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewStorageConnectorHandler()
	r.Post("/", h.CreateStorageConnector)
	r.Get("/", h.GetStorageConnectors)
	return r
}
