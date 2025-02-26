package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/phuslu/log"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/auth"
	"github.com/uploadpilot/manager/internal/utils"
)

type WorkspacePerm string

const (
	CanRead   WorkspacePerm = "read_ws"
	CanManage WorkspacePerm = "manage_ws"
	CanUpload WorkspacePerm = "upload_ws"
)

type AccountPerm string

const (
	CanReadAcc   AccountPerm = "read_acc"
	CanManageAcc AccountPerm = "manage_acc"
)

type Middlewares struct {
	authSvc *auth.Service
}

func NewAppMiddlewares(authSvc *auth.Service) *Middlewares {
	return &Middlewares{authSvc: authSvc}
}

func (m *Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", config.AllowedOrigins)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Header().Set("Access-Control-Allow-Headers", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
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

func (m *Middlewares) AccountAuthMiddleware(perms ...AccountPerm) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var claims *dto.UserClaims
			akc, err := m.verifyAPIKey(r)
			if err == nil && m.checkAccountAccess(akc, perms...) {
				claims = akc.UserClaims
			} else {
				jwtc, err := m.verifyJWTToken(r)
				if err != nil {
					utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
					return
				}
				claims = jwtc.UserClaims
			}

			next.ServeHTTP(w, r.WithContext(m.prepareContext(r, claims)))
		})
	}
}

func (m *Middlewares) WorkspaceAuthMiddleware(perms ...WorkspacePerm) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			workspaceID := chi.URLParam(r, "workspaceId")
			if workspaceID == "" {
				utils.HandleHttpError(w, r, http.StatusBadRequest, errors.New("workspaceId is required"))
				return
			}

			claims, err := m.authenticateWorkspacePerm(r, workspaceID, perms...)
			if err != nil {
				utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(m.prepareContext(r, claims)))
		})
	}
}

func (m *Middlewares) JWTOnlyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.verifyJWTToken(r)
		if err != nil {
			utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}
		next.ServeHTTP(w, r)
	})
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

func (m *Middlewares) verifyAPIKey(r *http.Request) (*dto.APIKeyClaims, error) {
	apiKey := r.Header.Get("X-Api-Key")
	claims, err := m.authSvc.ValidateAPIKey(r.Context(), apiKey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *Middlewares) verifyJWTToken(r *http.Request) (*dto.JWTClaims, error) {
	jwtToken := r.Header.Get("Authorization")
	claims, err := m.authSvc.ValidateJWTToken(jwtToken)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *Middlewares) authenticateWorkspacePerm(r *http.Request, workspaceID string, perms ...WorkspacePerm) (*dto.UserClaims, error) {
	akc, err := m.verifyAPIKey(r)
	if err == nil {
		for _, perm := range akc.Permissions {
			if perm.WorkspaceID == workspaceID && m.checkWorkspaceAccess(&perm, perms...) {
				return akc.UserClaims, nil
			}
		}
	}

	jwtc, err := m.verifyJWTToken(r)
	if err != nil {
		return nil, err
	}
	return jwtc.UserClaims, nil
}

func (m *Middlewares) prepareContext(r *http.Request, claims *dto.UserClaims) context.Context {
	ctx := context.WithValue(r.Context(), dto.UserIDContextKey, claims.UserID)
	ctx = context.WithValue(ctx, dto.EmailContextKey, claims.Email)
	ctx = context.WithValue(ctx, dto.NameContextKey, claims.Name)
	return ctx
}

func (m *Middlewares) checkAccountAccess(claims *dto.APIKeyClaims, access ...AccountPerm) bool {
	for _, a := range access {
		switch a {
		case CanReadAcc:
			if claims.CanReadAcc {
				return true
			}
		case CanManageAcc:
			if claims.CanManageAcc {
				return true
			}
		}
	}
	return false
}

func (m *Middlewares) checkWorkspaceAccess(perm *models.APIKeyPerm, access ...WorkspacePerm) bool {
	for _, a := range access {
		switch a {
		case CanRead:
			if perm.CanRead {
				return true
			}
		case CanManage:
			if perm.CanManage {
				return true
			}
		case CanUpload:
			if perm.CanUpload {
				return true
			}
		}
	}
	return false
}
