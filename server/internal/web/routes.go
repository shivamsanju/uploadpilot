package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/shivamsanju/uploader/internal/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/", GetAuthRoutes())
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Mount("/users", GetUserRoutes())
		r.Mount("/workflows", GetWorkflowRoutes())
		r.Mount("/storage/connectors", GetStorageConnectorRoutes())
		r.Mount("/storage/datastores", GetDatastoreRoutes())
		r.Mount("/importPolicies", GetImportPoliciesRoutes())
	})

	return r
}

func GetAuthRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewUserHandler()
	r.Post("/signup", h.Signup)
	r.Post("/login", h.Login)
	return r
}

func GetUserRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewUserHandler()
	r.Get("/me", h.GetUserDetails)
	return r
}

func GetWorkflowRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewWorkflowHandler()
	r.Post("/", h.CreateWorkflow)
	r.Get("/", h.GetAllWorkflows)
	r.Get("/{id}", h.GetWorkflowByID)
	r.Delete("/{id}", h.DeleteWorkflow)
	return r
}

func GetStorageConnectorRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewStorageConnectorHandler()
	r.Post("/", h.CreateStorageConnector)
	r.Get("/", h.GetAllStorageConnectors)
	r.Get("/{id}", h.GetStorageConnectorByID)
	r.Delete("/{id}", h.DeleteStorageConnector)
	return r
}

func GetDatastoreRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewDatastoreHandler()
	r.Post("/", h.CreateDatastore)
	r.Get("/", h.GetAllDatastores)
	r.Get("/{id}", h.GetDatastoreByID)
	r.Delete("/{id}", h.DeleteDatastore)
	return r
}

func GetImportPoliciesRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewImportPolicyHandler()
	r.Post("/", h.CreateImportPolicy)
	r.Get("/", h.GetImportPolicies)
	r.Get("/{id}", h.GetImportPolicy)
	r.Delete("/{id}", h.DeleteImportPolicy)
	return r
}
