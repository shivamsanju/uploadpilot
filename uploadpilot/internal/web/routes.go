package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shivamsanju/uploader/internal/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	ih := handlers.NewTusdHandler()
	r.Mount("/upload", http.StripPrefix(("/upload"), ih.GetTusHandler()))

	// Uploader details and signup login - no auth
	r.Mount("/auth", GetAuthRoutes())
	r.Mount("/uploaders", GetUploaderRoutes())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Mount("/users", GetUserRoutes())
		r.Mount("/uploaders/{uploaderId}/imports", GetImportRoutes())
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
	r.Get("/{uploaderId}", h.GetUploaderByID)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/", h.CreateUploader)
		r.Get("/", h.GetAllUploaders)
		r.Put("/{uploaderId}/config", h.UpdateUploaderConfig)
		r.Delete("/{uploaderId}", h.DeleteUploader)
		r.Get("/allowedSources", h.GetAllAllowedSources)
	})
	return r
}

func GetImportRoutes() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewImportHandler()
	r.Get("/", h.GetAllImportsForUploader)
	r.Get("/{importId}", h.GetImportDetailsByID)
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
