package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/web/handlers"
	"github.com/uploadpilot/core/web/middlewares"
	"github.com/uploadpilot/core/web/webutils"
)

func NewAppRoutesV1(services *services.Services, middlewares *middlewares.Middlewares) *chi.Mux {
	router := chi.NewRouter()

	// Handlers for uploads
	userHandler := handlers.NewUserHandler()
	tenantHandler := handlers.NewTenantHandler(services.TenantService)
	authHandler := handlers.NewAPIKeyHandler(services.APIKeyService)
	workspaceHandler := handlers.NewWorkspaceHandler(services.WorkspaceService)
	uploadHandler := handlers.NewUploadHandler(services.UploadService, services.WorkspaceService)
	procHandler := handlers.NewProcessorsHandler(services.ProcessorService)

	router.Use(supertokens.Middleware)
	router.Use(middlewares.CorsMiddleware)
	router.Use(middlewares.VerifySession)

	router.Route("/user", func(r chi.Router) {
		r.Get("/", webutils.CreateJSONHandler(userHandler.GetUserDetails))
		r.Put("/", webutils.CreateJSONHandler(userHandler.UpdateUserDetails))
		r.Put("/active-tenant", webutils.CreateJSONHandler(tenantHandler.SetActiveTenant))
	})

	router.Route("/tenants", func(r chi.Router) {
		r.Post("/", webutils.CreateJSONHandler(tenantHandler.OnboardTenant))
	})

	// Specific tenant routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.VerifyTenantAccess)

		r.Route("/tenants/{tenantId}", func(r chi.Router) {
			r.Route("/api-keys", func(r chi.Router) {
				r.Get("/", webutils.CreateJSONHandler(authHandler.GetAPIKeys))
				r.Post("/", webutils.CreateJSONHandler(authHandler.CreateAPIKey))
				r.Route("/{apiKeyId}", func(r chi.Router) {
					r.Post("/revoke", webutils.CreateJSONHandler(authHandler.RevokeAPIKey))
				})
			})

			r.Route("/workspaces", func(r chi.Router) {
				r.Get("/", webutils.CreateJSONHandler(workspaceHandler.GetAllWorkspaces))
				r.Post("/", webutils.CreateJSONHandler(workspaceHandler.CreateWorkspace))

				r.Route("/{workspaceId}", func(r chi.Router) {
					r.Get("/config", webutils.CreateJSONHandler(workspaceHandler.GetWorkspaceConfig))
					r.Put("/config", webutils.CreateJSONHandler(workspaceHandler.SetWorkspaceConfig))

					r.Route("/uploads", func(r chi.Router) {
						r.Get("/", webutils.CreateJSONHandler(uploadHandler.GetPaginatedUploads))
						r.Post("/", webutils.CreateJSONHandler(uploadHandler.CreateUpload))
						r.Post("/log", webutils.CreateJSONHandler(workspaceHandler.LogUpload))
						r.Route("/{uploadId}", func(r chi.Router) {
							r.Get("/", webutils.CreateJSONHandler(uploadHandler.GetUploadDetailsByID))
							r.Post("/finish", webutils.CreateJSONHandler(uploadHandler.FinishUpload))
							r.Get("/download", webutils.CreateJSONHandler(uploadHandler.GetUploadURL))
							r.Post("/process", webutils.CreateJSONHandler(uploadHandler.ProcessUpload))
						})
					})

					r.Route("/processors", func(r chi.Router) {
						r.Get("/", webutils.CreateJSONHandler(procHandler.GetProcessors))
						r.Post("/", webutils.CreateJSONHandler(procHandler.CreateProcessor))
						r.Get("/activities", webutils.CreateJSONHandler(procHandler.GetAllActivities))
						r.Get("/templates", webutils.CreateJSONHandler(procHandler.GetTemplates))
						r.Route("/{processorId}", func(r chi.Router) {
							r.Get("/", webutils.CreateJSONHandler(procHandler.GetProcessorDetailsByID))
							r.Put("/", webutils.CreateJSONHandler(procHandler.UpdateProcessor))
							r.Delete("/", webutils.CreateJSONHandler(procHandler.DeleteProcessor))
							r.Put("/enable", webutils.CreateJSONHandler(procHandler.EnableProcessor))
							r.Put("/disable", webutils.CreateJSONHandler(procHandler.DisableProcessor))
							r.Put("/workflow", webutils.CreateJSONHandler(procHandler.UpdateWorkflow))
							r.Get("/runs", webutils.CreateJSONHandler(procHandler.GetWorkflowRuns))
							r.Get("/logs", webutils.CreateJSONHandler(procHandler.GetWorkflowLogs))
						})
					})
				})
			})
		})
	})

	return router
}
