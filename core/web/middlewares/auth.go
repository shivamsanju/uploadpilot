package middlewares

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/web/webutils"
)

func (m *Middlewares) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if webutils.GetAPIKeyFromReq(r) != "" {
			m.VerifyAPIKeySession(next).ServeHTTP(w, r)
		} else {
			m.VerifyBearerTokenSession(next).ServeHTTP(w, r)
		}
	})
}

func (m *Middlewares) VerifyBearerTokenSession(next http.Handler) http.Handler {
	return session.VerifySession(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := m.getSessionFromBearerToken(r)
		if err != nil {
			log.Error().Msg(err.Error())
			webutils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.SessionCtxKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func (m *Middlewares) getSessionFromBearerToken(r *http.Request) (*dto.Session, error) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	userID := sessionContainer.GetUserID()
	if userID == "" {
		return nil, errors.New(msg.ErrUserInfoNotFoundInRequest)
	}

	email, err := webutils.GetEmailFromUserID(userID)
	if err != nil {
		return nil, err
	}

	metadata, err := usermetadata.GetUserMetadata(userID)
	if err != nil {
		return nil, err
	}

	return &dto.Session{
		Sub:      userID,
		UserID:   userID,
		Email:    email,
		Metadata: metadata,
	}, nil
}

func (m *Middlewares) VerifyAPIKeySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := webutils.GetAPIKeyFromReq(r)
		key, err := m.apiKeySvc.VerifyAPIKey(r.Context(), apiKey)
		if err != nil {
			log.Error().Msg(err.Error())
			webutils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}
		sess, err := m.getSessionToCtxFromAPIKey(key)
		if err != nil {
			log.Error().Msg(err.Error())
			webutils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.SessionCtxKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middlewares) getSessionToCtxFromAPIKey(key *models.APIKey) (*dto.Session, error) {

	email, err := webutils.GetEmailFromUserID(key.UserID)
	if err != nil {
		return nil, err
	}

	metadata, err := usermetadata.GetUserMetadata(key.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.Session{
		Sub:      key.ApiKeyHash,
		UserID:   key.UserID,
		Email:    email,
		Metadata: metadata,
	}, nil
}

func (m *Middlewares) VerifyTenantAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.verifyTenantHelper(r)
		if err != nil {
			webutils.HandleHttpError(w, r, http.StatusUnauthorized, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) verifyTenantHelper(r *http.Request) error {
	reqTenantID := webutils.GetTenantIDFromReq(r)
	if reqTenantID == "" {
		return errors.New(msg.ErrTenantIDNotFoundInRequest)
	}

	sess, err := webutils.GetSessionFromCtx(r.Context())
	if err != nil {
		return err
	}

	tenants, err := webutils.GetUserTenantsFromMetadata(sess.Metadata)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user tenants")
		return err
	}

	_, hasAcess := tenants[reqTenantID]
	if !hasAcess {
		return errors.New(msg.ErrInvalidTenantIDInRequest)
	}

	if wID, ok := isWorkspaceRoute(r.URL.Path); ok {
		log.Debug().Str("path", r.URL.Path).Str("WorkspaceID", wID).Msg("is workspace route")
		_, err := uuid.Parse(wID)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse workspace id")
			return errors.New(msg.ErrAccessDenied)
		}
		return m.verifyWorkspacePartOfTenant(r.Context(), reqTenantID, wID)
	}

	return nil
}

func (m *Middlewares) verifyWorkspacePartOfTenant(ctx context.Context, tenantID, workspaceID string) error {
	tID, err := m.workspaceSvc.GetTenantID(ctx, workspaceID)
	if err != nil {
		return err
	}

	log.Debug().Interface("tID", tID).Str("tenantID", tenantID).Msg("tID")

	if tID != tenantID {
		return errors.New(msg.ErrAccessDenied)
	}

	return nil
}

func isWorkspaceRoute(path string) (string, bool) {
	pattern := `^/tenants/([^/]+)/workspaces/([^/]+)/`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(path)
	if len(matches) > 2 {
		return matches[2], true
	}
	return "", false
}
