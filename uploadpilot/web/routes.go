package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/uploadpilot/web/handlers"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

	// Handlers for uploads
	ih := handlers.NewTusdHandler()
	authHandler := handlers.NewAuthHandler()
	uploaderHandler := handlers.NewuploaderHandler()
	importHandler := handlers.NewImportHandler()
	storageHandler := handlers.NewStorageConnectorHandler()

	// Public routes
	r.Group(func(r chi.Router) {
		r.Mount("/upload", http.StripPrefix("/upload", ih.GetTusHandler()))

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Get("/{provider}/authorize", authHandler.Login)
			r.Get("/{provider}/callback", authHandler.HandleCallback)
			r.Get("/logout", authHandler.Logout)
			r.Get("/logout/{provider}", authHandler.LogoutProvider)
		})

		// Uploader details
		r.Get("/uploaders/{uploaderId}", uploaderHandler.GetUploaderByID)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		// Session routes
		r.Get("/session", authHandler.GetSession)

		// Uploader routes
		r.Route("/uploaders", func(r chi.Router) {
			r.Post("/", uploaderHandler.CreateUploader)
			r.Get("/", uploaderHandler.GetAllUploaders)
			r.Get("/sources/allowed", uploaderHandler.GetAllAllowedSources)
			r.Put("/{uploaderId}/config", uploaderHandler.UpdateUploaderConfig)
			r.Delete("/{uploaderId}", uploaderHandler.DeleteUploader)
		})

		// Import routes
		r.Route("/uploaders/{uploaderId}/imports", func(r chi.Router) {
			r.Get("/", importHandler.GetAllImportsForUploader)
			r.Get("/{importId}", importHandler.GetImportDetailsByID)
		})

		// Storage connector routes
		r.Route("/storageConnectors", func(r chi.Router) {
			r.Post("/", storageHandler.CreateStorageConnector)
			r.Get("/", storageHandler.GetAllStorageConnectors)
			r.Get("/{id}", storageHandler.GetStorageConnectorByID)
			r.Delete("/{id}", storageHandler.DeleteStorageConnector)
		})
	})

	return r
}
