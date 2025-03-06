package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/phuslu/log"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/web/webutils"
)

func (m *Middlewares) VerifySession(next http.Handler) http.Handler {
	return session.VerifySession(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.addSessionToCtx(r)
		if err != nil {
			webutils.HandleHttpError(w, r, http.StatusUnauthorized, err)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
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

func (m *Middlewares) addSessionToCtx(r *http.Request) (context.Context, error) {
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

	sess := &dto.Session{
		UserID:   userID,
		Email:    email,
		Metadata: metadata,
	}

	ctx := context.WithValue(r.Context(), dto.SessionCtxKey, sess)
	return ctx, nil
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

	return nil
}
