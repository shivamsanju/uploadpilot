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

			r.Get("/", middlewares.CheckPermissions(CreateJSONHandler(authHandler.GetAPIKeys), TenantReadAccess))
			r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(authHandler.CreateAPIKey), TenantReadAccess))
			r.Route("/{apiKeyId}", func(r chi.Router) {
				r.Post("/revoke", middlewares.CheckPermissions(CreateJSONHandler(authHandler.RevokeAPIKey), TenantReadAccess))
			})
		})

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {

			r.Get("/", middlewares.CheckPermissions(workspaceHandler.GetAllWorkspaces, TenantReadAccess))
			r.Post("/", middlewares.CheckPermissions(workspaceHandler.CreateWorkspace, TenantReadAccess))

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {

				r.Get("/config", middlewares.CheckPermissions(workspaceHandler.GetWorkspaceConfig, WorkspaceUploadAccess))
				r.Put("/config", middlewares.CheckPermissions(workspaceHandler.SetWorkspaceConfig, TenantReadAccess))
				r.Get("/allowedSources", middlewares.CheckPermissions(CreateJSONHandler(workspaceHandler.GetAllAllowedSources), WorkspaceUploadAccess))

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Get("/", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.GetPaginatedUploads), TenantReadAccess))
					r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.CreateUpload), WorkspaceUploadAccess))
					r.Post("/log", CreateJSONHandler(workspaceHandler.LogUpload))

					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", middlewares.CheckPermissions(uploadHandler.GetUploadDetailsByID, TenantReadAccess))
						r.Post("/finish", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.FinishUpload), WorkspaceUploadAccess))
						r.Get("/download", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.GetUploadURL), TenantReadAccess))
						r.Post("/process", middlewares.CheckPermissions(CreateJSONHandler(uploadHandler.ProcessUpload), TenantReadAccess))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", middlewares.CheckPermissions(procHandler.GetProcessors, TenantReadAccess))
					r.Post("/", middlewares.CheckPermissions(CreateJSONHandler(procHandler.CreateProcessor), TenantReadAccess))
					r.Get("/activities", middlewares.CheckPermissions(procHandler.GetAllActivities, TenantReadAccess))
					r.Get("/templates", middlewares.CheckPermissions(procHandler.GetTemplates, TenantReadAccess))
					r.Route("/{processorId}", func(r chi.Router) {
						r.Get("/", middlewares.CheckPermissions(procHandler.GetProcessorDetailsByID, TenantReadAccess))
						r.Put("/", middlewares.CheckPermissions(procHandler.UpdateProcessor, TenantReadAccess))
						r.Delete("/", middlewares.CheckPermissions(procHandler.DeleteProcessor, TenantReadAccess))
						r.Put("/enable", middlewares.CheckPermissions(procHandler.EnableProcessor, TenantReadAccess))
						r.Put("/disable", middlewares.CheckPermissions(procHandler.DisableProcessor, TenantReadAccess))
						r.Put("/workflow", middlewares.CheckPermissions(procHandler.UpdateWorkflow, TenantReadAccess))
						r.Get("/runs", middlewares.CheckPermissions(procHandler.GetWorkflowRuns, TenantReadAccess))
						r.Get("/logs", middlewares.CheckPermissions(procHandler.GetWorkflowLogs, TenantReadAccess))
					})
				})
			})
		})
	})

	return router
}
