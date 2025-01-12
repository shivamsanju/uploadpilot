package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shivamsanju/uploader/internal/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/", GetAuthRoutes())
	ih := handlers.NewImportHandler()
	r.Mount("/imports", http.StripPrefix(("/imports"), ih.GetTusHandler()))
	// r.Mount("/imports/{*path}", http.StripPrefix(("/imports/"), ih.GetTusHandler()))
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Mount("/users", GetUserRoutes())
		r.Mount("/uploaders", GetUploaderRoutes())
		r.Mount("/storage/connectors", GetStorageConnectorRoutes())
		r.Mount("/storage/datastores", GetDatastoreRoutes())
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

func GetUploaderRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewuploaderHandler()
	r.Post("/", h.CreateUploader)
	r.Get("/", h.GetAllUploaders)
	r.Get("/{id}", h.GetUploaderByID)
	r.Delete("/{id}", h.DeleteUploader)
	r.Get("/allowedSources", h.GetAllAllowedSources)
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
