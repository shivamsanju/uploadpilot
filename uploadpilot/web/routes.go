package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/uploadpilot/web/handlers"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(AllowAllCorsMiddleware)

	// Handlers for uploads
	ih := handlers.NewTusdHandler()
	authHandler := handlers.NewAuthHandler()
	workspaceHandler := handlers.NewWorkspaceHandler()
	importHandler := handlers.NewImportHandler()

	// Public routes
	router.Group(func(r chi.Router) {
		r.Mount("/upload", http.StripPrefix("/upload", ih.GetTusHandler()))
		// Uploader details
		r.Get("/workspaces/{workspaceId}/config", workspaceHandler.GetUploaderConfig)
	})

	// Auth routes
	router.Group(func(r chi.Router) {
		r.Use(CorsMiddleware)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/{provider}/authorize", authHandler.Login)
			r.Get("/{provider}/callback", authHandler.HandleCallback)
			r.Get("/logout", authHandler.Logout)
			r.Get("/logout/{provider}", authHandler.LogoutProvider)
		})
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(CorsMiddleware)
		r.Use(AuthMiddleware)

		// Session
		r.Get("/session", authHandler.GetSession)

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {
			r.Post("/", workspaceHandler.CreateWorkspace)
			r.Get("/", workspaceHandler.GetWorkspacesForUser)

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {
				r.Put("/config", workspaceHandler.UpdateUploaderConfig)
				r.Get("/allowedSources", workspaceHandler.GetAllAllowedSources)

				// Users
				r.Route("/users", func(r chi.Router) {
					r.Get("/", workspaceHandler.GetAllUsersInWorkspace)
					r.Post("/", workspaceHandler.AddUserToWorkspace)
					r.Delete("/", workspaceHandler.RemoveUserFromWorkspace)
				})

				// Imports
				r.Route("/imports", func(r chi.Router) {
					r.Get("/", importHandler.GetAllImportsForWorkspace)
					r.Route("/{importId}", func(r chi.Router) {
						r.Get("/", importHandler.GetImportDetailsByID)
					})
				})
			})
		})
	})

	return router
}
