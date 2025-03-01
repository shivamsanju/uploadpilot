package web

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/phuslu/log"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/manager/internal/svc/apikey"
	"github.com/uploadpilot/manager/internal/utils"
)

type Middlewares struct {
	apiKeySvc *apikey.Service
}

func NewAppMiddlewares(apiKeySvc *apikey.Service) *Middlewares {
	return &Middlewares{
		apiKeySvc: apiKeySvc,
	}
}

func (m *Middlewares) RecoveryMiddleware(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func (m *Middlewares) RequestIDMiddleware(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func (m *Middlewares) RequestTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

func (m *Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Validate and set allowed origins
		if origin != "" && slices.Contains(config.AllowedOrigins, origin) {
			response.Header().Set("Access-Control-Allow-Origin", origin) // Dynamically set origin
			response.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", strings.Join(
			append([]string{"Content-Type", "X-Api-Key", "X-Tenant-Id"}, supertokens.GetAllCORSHeaders()...),
			",",
		))

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			response.WriteHeader(http.StatusNoContent) // No content needed for preflight
			return                                     // Return early to avoid calling next handler
		}

		next.ServeHTTP(response, r)
	})
}

func (m *Middlewares) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			if ww.Status() >= 400 {
				log.Error().
					Str("request_id", middleware.GetReqID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("protocol", r.Proto).
					Str("remote_addr", r.RemoteAddr).
					Str("status", utils.GetStatusLabel(ww.Status())).
					Int("bytes_written", int(ww.BytesWritten())).
					Int64("time_taken", time.Since(t1).Milliseconds()).
					Msg("request failed")
			} else {
				log.Info().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("bytes_written", int(ww.BytesWritten())).
					Int64("time_taken", time.Since(t1).Milliseconds()).
					Msg("request completed")
			}
		}()
		next.ServeHTTP(ww, r)
	})
}

func (m *Middlewares) CheckPermissions(h http.HandlerFunc, perm APIPermission) http.HandlerFunc {
	return m.checkPermissionsHelper(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.addSessionToCtx(r)
		if err != nil {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}
		h(w, r.WithContext(ctx))
	}, perm)
}

func (m *Middlewares) checkPermissionsHelper(h http.HandlerFunc, perm APIPermission) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, authType := range perm.AllowedAuthTypes {
			switch authType {
			case APIAuthTypeAPIKey:
				apiKey := utils.GetAPIKeyFromReq(r)
				if apiKey == "" {
					continue
				}
				// If we have an api key, it should be correct
				workspaceID := chi.URLParam(r, "workspaceId")
				tenantID := r.Header.Get("X-Tenant-Id")
				for i, _ := range perm.Permissions {
					if perm.Permissions[i].ResouceID == "<workspaceId>" {
						perm.Permissions[i].ResouceID = workspaceID
					}
					if perm.Permissions[i].ResouceID == "<tenantId>" {
						perm.Permissions[i].ResouceID = tenantID
					}
				}
				log.Debug().Interface("requested_permissions", perm.Permissions).Msg("verifying api key permissions")
				apiKeyDetails, err := m.apiKeySvc.ValidateAPIKey(r.Context(), apiKey, perm.Permissions...)
				if err != nil {
					log.Error().Err(err).Msg("failed to validate api key")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				r.Header.Set("X-Api-Key-UserId", apiKeyDetails.UserID)
				r.Header.Set("X-Tenant-Id", apiKeyDetails.TenantID)
				h(w, r)
				return
			case APIAuthTypeBearer:
				session.VerifySession(nil, h)(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

}

func (m *Middlewares) addSessionToCtx(r *http.Request) (context.Context, error) {
	var sess dto.Session

	if r.Header.Get("X-Api-Key") != "" {
		if r.Header.Get("X-Api-Key-UserId") != "" {
			sess.UserID = r.Header.Get("X-Api-Key-UserId")
		} else {
			return nil, errors.New(msg.ErrUserInfoNotFoundInRequest)
		}
	} else {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()
		if userID == "" {
			return nil, errors.New(msg.ErrUserInfoNotFoundInRequest)
		}
		sess.UserID = userID
	}

	tpusr, err := thirdparty.GetUserByID(sess.UserID)
	if err == nil && tpusr != nil {
		sess.Email = tpusr.Email
	} else {
		epusr, err := emailpassword.GetUserByID(sess.UserID)
		if err == nil && epusr != nil {
			sess.Email = epusr.Email
		} else {
			return nil, errors.New(msg.ErrUserInfoNotFoundInRequest)
		}
	}

	metadata, err := usermetadata.GetUserMetadata(sess.UserID)
	if err != nil {
		return nil, err
	}

	reqTenantID := utils.GetTenantIDFromReq(r)
	if reqTenantID == "" {
		return nil, errors.New(msg.ErrTenantIDNotFoundInRequest)
	}
	tenants, err := utils.GetUserTenantsFromMetadata(metadata)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user tenants")
		return nil, err
	}

	tenant, hasAcess := tenants[reqTenantID]
	if !hasAcess {
		return nil, errors.New(msg.ErrInvalidTenantIDInRequest)
	}

	sess.TenantID = tenant.ID
	sess.Metadata = metadata

	ctx := context.WithValue(r.Context(), dto.SessionCtxKey, sess)

	return ctx, nil
}
