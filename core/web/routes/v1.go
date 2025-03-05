package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/supertokens/supertokens-golang/recipe/session"
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

	router.Group(func(r chi.Router) {
		// User
		r.Route("/user", func(r chi.Router) {
			r.Get("/", session.VerifySession(nil, webutils.CreateJSONHandler(userHandler.GetUserDetails)))
			r.Put("/", session.VerifySession(nil, webutils.CreateJSONHandler(userHandler.UpdateUserDetails)))

		})

		// Tenants
		r.Route("/tenants", func(r chi.Router) {
			r.Post("/", session.VerifySession(nil, webutils.CreateJSONHandler(tenantHandler.OnboardTenant)))
			r.Put("/active", session.VerifySession(nil, webutils.CreateJSONHandler(tenantHandler.SetActiveTenant)))
		})

	})

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.CorsMiddleware)
		r.Use(middlewares.AuthMiddleware)

		// ApiKeys
		r.Route("/api-keys", func(r chi.Router) {

			r.Get("/", webutils.CreateJSONHandler(authHandler.GetAPIKeys))
			r.Post("/", webutils.CreateJSONHandler(authHandler.CreateAPIKey))
			r.Route("/{apiKeyId}", func(r chi.Router) {
				r.Post("/revoke", webutils.CreateJSONHandler(authHandler.RevokeAPIKey))
			})
		})

		// Workspaces
		r.Route("/workspaces", func(r chi.Router) {
			r.Get("/", workspaceHandler.GetAllWorkspaces)
			r.Post("/", workspaceHandler.CreateWorkspace)

			// Single workspace
			r.Route("/{workspaceId}", func(r chi.Router) {

				r.Get("/config", workspaceHandler.GetWorkspaceConfig)
				r.Put("/config", workspaceHandler.SetWorkspaceConfig)
				r.Get("/allowedSources", webutils.CreateJSONHandler(workspaceHandler.GetAllAllowedSources))

				// Uploads
				r.Route("/uploads", func(r chi.Router) {
					r.Get("/", webutils.CreateJSONHandler(uploadHandler.GetPaginatedUploads))
					r.Post("/", webutils.CreateJSONHandler(uploadHandler.CreateUpload))
					r.Post("/log", webutils.CreateJSONHandler(workspaceHandler.LogUpload))

					r.Route("/{uploadId}", func(r chi.Router) {
						r.Get("/", uploadHandler.GetUploadDetailsByID)
						r.Post("/finish", webutils.CreateJSONHandler(uploadHandler.FinishUpload))
						r.Get("/download", webutils.CreateJSONHandler(uploadHandler.GetUploadURL))
						r.Post("/process", webutils.CreateJSONHandler(uploadHandler.ProcessUpload))
					})
				})

				// processors
				r.Route("/processors", func(r chi.Router) {
					r.Get("/", procHandler.GetProcessors)
					r.Post("/", webutils.CreateJSONHandler(procHandler.CreateProcessor))
					r.Get("/activities", procHandler.GetAllActivities)
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
