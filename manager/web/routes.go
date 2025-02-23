package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/manager/web/handlers"
)

func Routes(services *svc.Services) *chi.Mux {
	router := chi.NewRouter()
	router.Use(CorsMiddleware)

	// Handlers for uploads
	authHandler := handlers.NewAuthHandler(services.UserService)
	workspaceHandler := handlers.NewWorkspaceHandler(services.WorkspaceService)
	uploadHandler := handlers.NewUploadHandler(services.UploadService)
	procHandler := handlers.NewProcessorsHandler(services.ProcessorService)

	// Auth routes
	router.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/{provider}/authorize", authHandler.Login)
			r.Get("/{provider}/callback", authHandler.HandleCallback)
			r.Get("/logout", authHandler.Logout)
			r.Get("/logout/{provider}", authHandler.LogoutProvider)
		})
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		// Session
		r.Get("/session", authHandler.GetSession)

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {
			r.Post("/", workspaceHandler.CreateWorkspace)
			r.Get("/", workspaceHandler.GetWorkspacesForUser)

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {
				r.Get("/config", workspaceHandler.GetUploaderConfig)
				r.Put("/config", workspaceHandler.UpdateUploaderConfig)
				r.Get("/allowedSources", CreateJSONHandler(workspaceHandler.GetAllAllowedSources))
				r.Get("/subscription", CreateJSONHandler(workspaceHandler.GetSubscription))
				r.Post("/log-upload", CreateJSONHandler(workspaceHandler.LogUpload))

				// Users
				r.Route("/users", func(r chi.Router) {
					r.Get("/", workspaceHandler.GetAllUsersInWorkspace)
					r.Post("/", workspaceHandler.AddUserToWorkspace)
					r.Route("/{userId}", func(r chi.Router) {
						r.Put("/", workspaceHandler.ChangeUserRoleInWorkspace)
						r.Delete("/", workspaceHandler.RemoveUserFromWorkspace)
					})
				})

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Post("/", CreateJSONHandler(uploadHandler.CreateUpload))
					r.Get("/", CreateJSONHandler(uploadHandler.GetPaginatedUploads))
					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", uploadHandler.GetUploadDetailsByID)
						r.Post("/finish", CreateJSONHandler(uploadHandler.FinishUpload))
						r.Get("/download", CreateJSONHandler(uploadHandler.GetUploadURL))
						r.Post("/process", CreateJSONHandler(uploadHandler.ProcessUpload))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", procHandler.GetProcessors)
					r.Post("/", CreateJSONHandler(procHandler.CreateProcessor))
					r.Get("/tasks", procHandler.GetAllTasks)
					r.Get("/templates", procHandler.GetTemplates)
					r.Route("/{processorId}", func(r chi.Router) {
						r.Get("/", procHandler.GetProcessorDetailsByID)
						r.Put("/", procHandler.UpdateProcessor)
						r.Delete("/", procHandler.DeleteProcessor)
						r.Put("/enable", procHandler.EnableProcessor)
						r.Put("/disable", procHandler.DisableProcessor)
						r.Put("/workflow", procHandler.UpdateWorkflow)
						r.Get("/runs", procHandler.GetWorkflowRuns)
						r.Get("/logs", procHandler.GetWorkflowLogs)
					})
				})
			})
		})
	})

	return router
}
