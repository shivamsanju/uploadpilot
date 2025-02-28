package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/manager/web/handlers"
)

func Routes(services *svc.Services, middlewares *Middlewares) *chi.Mux {
	router := chi.NewRouter()

	// Handlers for uploads
	userHandler := handlers.NewUserHandler()
	tenantHandler := handlers.NewTenantHandler(services.TenantService)
	authHandler := handlers.NewAPIKeyHandler(services.APIKeyService)
	workspaceHandler := handlers.NewWorkspaceHandler(services.WorkspaceService)
	uploadHandler := handlers.NewUploadHandler(services.UploadService, services.WorkspaceService)
	procHandler := handlers.NewProcessorsHandler(services.ProcessorService)

	router.Use(supertokens.Middleware)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.CorsMiddleware)

		// User
		r.Route("/user", func(r chi.Router) {
			r.Get("/", session.VerifySession(nil, CreateJSONHandler(userHandler.GetUserDetails)))
			r.Put("/", session.VerifySession(nil, CreateJSONHandler(userHandler.UpdateUserDetails)))

		})

		// Tenants
		r.Route("/tenants", func(r chi.Router) {
			r.Post("/", session.VerifySession(nil, CreateJSONHandler(tenantHandler.OnboardTenant)))
			r.Put("/active", session.VerifySession(nil, CreateJSONHandler(tenantHandler.SetActiveTenant)))
		})

		// ApiKeys
		r.Route("/api-keys", func(r chi.Router) {

			r.Get("/", middlewares.CheckPermissions(CreateJSONHandler(authHandler.GetAPIKeys), BearerTenantReadAccess))
			r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(authHandler.CreateAPIKey), BearerTenantReadAccess))
			r.Route("/{apiKeyId}", func(r chi.Router) {
				r.Post("/revoke", middlewares.CheckPermissions(CreateJSONHandler(authHandler.RevokeAPIKey), BearerTenantReadAccess))
			})
		})

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {

			r.Get("/", middlewares.CheckPermissions(workspaceHandler.GetAllWorkspaces, BearerTenantReadAccess))
			r.Post("/", middlewares.CheckPermissions(workspaceHandler.CreateWorkspace, BearerTenantReadAccess))

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {

				r.Get("/config", middlewares.CheckPermissions(workspaceHandler.GetWorkspaceConfig, BearerTenantReadAccess))
				r.Put("/config", middlewares.CheckPermissions(workspaceHandler.SetWorkspaceConfig, BearerTenantReadAccess))
				r.Get("/allowedSources", middlewares.CheckPermissions(CreateJSONHandler(workspaceHandler.GetAllAllowedSources), BearerTenantReadAccess))

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Get("/", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.GetPaginatedUploads), BearerTenantReadAccess))
					r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.CreateUpload), BearerTenantReadAccess))
					r.Post("/log", CreateJSONHandler(workspaceHandler.LogUpload))

					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", middlewares.CheckPermissions(uploadHandler.GetUploadDetailsByID, BearerTenantReadAccess))
						r.Post("/finish", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.FinishUpload), BearerTenantReadAccess))
						r.Get("/download", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.GetUploadURL), BearerTenantReadAccess))
						r.Post("/process", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.ProcessUpload), BearerTenantReadAccess))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", middlewares.CheckPermissions(procHandler.GetProcessors, BearerTenantReadAccess))
					r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(procHandler.CreateProcessor), BearerTenantReadAccess))
					r.Get("/tasks", middlewares.CheckPermissions(procHandler.GetAllTasks, BearerTenantReadAccess))
					r.Get("/templates", middlewares.CheckPermissions(procHandler.GetTemplates, BearerTenantReadAccess))
					r.Route("/{processorId}", func(r chi.Router) {
						r.Get("/", middlewares.CheckPermissions(procHandler.GetProcessorDetailsByID, BearerTenantReadAccess))
						r.Put("/", middlewares.CheckPermissions(procHandler.UpdateProcessor, BearerTenantReadAccess))
						r.Delete("/", middlewares.CheckPermissions(procHandler.DeleteProcessor, BearerTenantReadAccess))
						r.Put("/enable", middlewares.CheckPermissions(procHandler.EnableProcessor, BearerTenantReadAccess))
						r.Put("/disable", middlewares.CheckPermissions(procHandler.DisableProcessor, BearerTenantReadAccess))
						r.Put("/workflow", middlewares.CheckPermissions(procHandler.UpdateWorkflow, BearerTenantReadAccess))
						r.Get("/runs", middlewares.CheckPermissions(procHandler.GetWorkflowRuns, BearerTenantReadAccess))
						r.Get("/logs", middlewares.CheckPermissions(procHandler.GetWorkflowLogs, BearerTenantReadAccess))
					})
				})
			})
		})
	})

	return router
}
