package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/uploadpilot/manager/internal/svc"
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
		r.Get("/session", WithMiddleware(authHandler.GetSession, middlewares.AccountAuthMiddleware(CanReadAcc, CanManageAcc)))

		// ApiKeys
		r.Route("/apikeys", func(r chi.Router) {
			r.Get("/", WithMiddleware(CreateJSONHandler(authHandler.GetAPIKeys), middlewares.JWTOnlyAuthMiddleware))
			r.Post("/", WithMiddleware(CreateJSONHandler(authHandler.CreateAPIKey), middlewares.JWTOnlyAuthMiddleware))
			r.Route("/{apiKeyId}", func(r chi.Router) {
				r.Post("/revoke", WithMiddleware(CreateJSONHandler(authHandler.RevokeAPIKey), middlewares.JWTOnlyAuthMiddleware))
			})
		})

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {
			r.Get("/", WithMiddleware(workspaceHandler.GetAllWorkspaces, middlewares.AccountAuthMiddleware(CanReadAcc, CanManageAcc)))
			r.Post("/", WithMiddleware(workspaceHandler.CreateWorkspace, middlewares.AccountAuthMiddleware(CanManageAcc)))

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {
				r.Get("/config", WithMiddleware(workspaceHandler.GetUploaderConfig, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage, CanUpload)))
				r.Put("/config", WithMiddleware(workspaceHandler.UpdateUploaderConfig, middlewares.WorkspaceAuthMiddleware(CanManage)))
				r.Get("/allowedSources", WithMiddleware(CreateJSONHandler(workspaceHandler.GetAllAllowedSources), middlewares.WorkspaceAuthMiddleware(CanRead, CanManage, CanUpload)))
				r.Get("/subscription", WithMiddleware(CreateJSONHandler(workspaceHandler.GetSubscription), middlewares.WorkspaceAuthMiddleware(CanRead, CanManage, CanUpload)))

				// Users
				r.Route("/users", func(r chi.Router) {
					r.Get("/", WithMiddleware(workspaceHandler.GetAllUsersInWorkspace, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					r.Post("/", WithMiddleware(workspaceHandler.AddUserToWorkspace, middlewares.WorkspaceAuthMiddleware(CanManage)))
					r.Route("/{userId}", func(r chi.Router) {
						r.Put("/", WithMiddleware(workspaceHandler.EditUser, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Delete("/", WithMiddleware(workspaceHandler.RemoveUserFromWorkspace, middlewares.WorkspaceAuthMiddleware(CanManage)))
					})
				})

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Get("/", WithMiddleware(CreateJSONHandler(uploadHandler.GetPaginatedUploads), middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					r.Post("/", WithMiddleware(CreateJSONHandler(uploadHandler.CreateUpload), middlewares.WorkspaceAuthMiddleware(CanUpload)))
					r.Post("/log", CreateJSONHandler(workspaceHandler.LogUpload))

					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", WithMiddleware(uploadHandler.GetUploadDetailsByID, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
						r.Post("/finish", WithMiddleware(CreateJSONHandler(uploadHandler.FinishUpload), middlewares.WorkspaceAuthMiddleware(CanUpload, CanManage)))
						r.Get("/download", WithMiddleware(CreateJSONHandler(uploadHandler.GetUploadURL), middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
						r.Post("/process", WithMiddleware(CreateJSONHandler(uploadHandler.ProcessUpload), middlewares.WorkspaceAuthMiddleware(CanManage)))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", WithMiddleware(procHandler.GetProcessors, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					r.Post("/", WithMiddleware(CreateJSONHandler(procHandler.CreateProcessor), middlewares.WorkspaceAuthMiddleware(CanManage)))
					r.Get("/tasks", WithMiddleware(procHandler.GetAllTasks, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					r.Get("/templates", WithMiddleware(procHandler.GetTemplates, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					r.Route("/{processorId}", func(r chi.Router) {
						r.Get("/", WithMiddleware(procHandler.GetProcessorDetailsByID, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
						r.Put("/", WithMiddleware(procHandler.UpdateProcessor, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Delete("/", WithMiddleware(procHandler.DeleteProcessor, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Put("/enable", WithMiddleware(procHandler.EnableProcessor, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Put("/disable", WithMiddleware(procHandler.DisableProcessor, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Put("/workflow", WithMiddleware(procHandler.UpdateWorkflow, middlewares.WorkspaceAuthMiddleware(CanManage)))
						r.Get("/runs", WithMiddleware(procHandler.GetWorkflowRuns, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
						r.Get("/logs", WithMiddleware(procHandler.GetWorkflowLogs, middlewares.WorkspaceAuthMiddleware(CanRead, CanManage)))
					})
				})
			})
		})
	})

	return router
}
