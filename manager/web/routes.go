package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/manager/internal/svc/auth"
	"github.com/uploadpilot/manager/web/handlers"
)

func Routes(services *svc.Services, middlewares *Middlewares) *chi.Mux {
	router := chi.NewRouter()

	// Handlers for uploads
	authHandler := handlers.NewAuthHandler(services.UserService, services.AuthService)
	workspaceHandler := handlers.NewWorkspaceHandler(services.WorkspaceService, services.AuthService)
	uploadHandler := handlers.NewUploadHandler(services.UploadService, services.WorkspaceService, services.AuthService)
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
		// Session
		r.Get("/session", WithMiddleware(authHandler.GetSession, middlewares.AccountAuthMiddleware(auth.CanReadAcc, auth.CanManageAcc)))

		// ApiKeys
		r.Route("/api-keys", func(r chi.Router) {
			r.Get("/", WithMiddleware(CreateJSONHandler(authHandler.GetAPIKeys), middlewares.JWTOnlyAuthMiddleware))
			r.Post("/", WithMiddleware(CreateJSONHandler(authHandler.CreateAPIKey), middlewares.JWTOnlyAuthMiddleware))
			r.Route("/{apiKeyId}", func(r chi.Router) {
				r.Post("/revoke", WithMiddleware(CreateJSONHandler(authHandler.RevokeAPIKey), middlewares.JWTOnlyAuthMiddleware))
			})
		})

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {
			r.Get("/", WithMiddleware(workspaceHandler.GetAllWorkspaces, middlewares.AccountAuthMiddleware(auth.CanReadAcc, auth.CanManageAcc)))
			r.Post("/", WithMiddleware(workspaceHandler.CreateWorkspace, middlewares.AccountAuthMiddleware(auth.CanManageAcc)))

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {
				r.Get("/config", WithMiddleware(workspaceHandler.GetUploaderConfig, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage, auth.CanUpload)))
				r.Put("/config", WithMiddleware(workspaceHandler.UpdateUploaderConfig, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
				r.Get("/allowedSources", WithMiddleware(CreateJSONHandler(workspaceHandler.GetAllAllowedSources), middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage, auth.CanUpload)))
				r.Get("/subscription", WithMiddleware(CreateJSONHandler(workspaceHandler.GetSubscription), middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage, auth.CanUpload)))

				// Users
				r.Route("/users", func(r chi.Router) {
					r.Get("/", WithMiddleware(workspaceHandler.GetAllUsersInWorkspace, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					r.Post("/", WithMiddleware(workspaceHandler.AddUserToWorkspace, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
					r.Route("/{userId}", func(r chi.Router) {
						r.Put("/", WithMiddleware(workspaceHandler.EditUser, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Delete("/", WithMiddleware(workspaceHandler.RemoveUserFromWorkspace, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
					})
				})

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Get("/", WithMiddleware(CreateJSONHandler(uploadHandler.GetPaginatedUploads), middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					r.Post("/", WithMiddleware(CreateJSONHandler(uploadHandler.CreateUpload), middlewares.WorkspaceAuthMiddleware(auth.CanUpload)))
					r.Post("/log", CreateJSONHandler(workspaceHandler.LogUpload))

					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", WithMiddleware(uploadHandler.GetUploadDetailsByID, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
						r.Post("/finish", WithMiddleware(CreateJSONHandler(uploadHandler.FinishUpload), middlewares.WorkspaceAuthMiddleware(auth.CanUpload, auth.CanManage)))
						r.Get("/download", WithMiddleware(CreateJSONHandler(uploadHandler.GetUploadURL), middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
						r.Post("/process", WithMiddleware(CreateJSONHandler(uploadHandler.ProcessUpload), middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", WithMiddleware(procHandler.GetProcessors, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					r.Post("/", WithMiddleware(CreateJSONHandler(procHandler.CreateProcessor), middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
					r.Get("/tasks", WithMiddleware(procHandler.GetAllTasks, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					r.Get("/templates", WithMiddleware(procHandler.GetTemplates, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					r.Route("/{processorId}", func(r chi.Router) {
						r.Get("/", WithMiddleware(procHandler.GetProcessorDetailsByID, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
						r.Put("/", WithMiddleware(procHandler.UpdateProcessor, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Delete("/", WithMiddleware(procHandler.DeleteProcessor, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Put("/enable", WithMiddleware(procHandler.EnableProcessor, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Put("/disable", WithMiddleware(procHandler.DisableProcessor, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Put("/workflow", WithMiddleware(procHandler.UpdateWorkflow, middlewares.WorkspaceAuthMiddleware(auth.CanManage)))
						r.Get("/runs", WithMiddleware(procHandler.GetWorkflowRuns, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
						r.Get("/logs", WithMiddleware(procHandler.GetWorkflowLogs, middlewares.WorkspaceAuthMiddleware(auth.CanRead, auth.CanManage)))
					})
				})
			})
		})
	})

	return router
}
