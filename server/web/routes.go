package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/shivamsanju/uploader/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/auth", getAuthRoutes())
	r.Mount("/workflows", getCodebaseRoutes())
	return r
}

func getCodebaseRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewWorkflowHandler()
	r.Post("/", h.CreateWorkflow)
	r.Get("/", h.ListWorkflows)
	r.Get("/{id}", h.GetWorkflow)
	return r
}

func getAuthRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewAuthHandler()
	r.Get("/userinfo", h.GetUserInfo)
	return r
}
