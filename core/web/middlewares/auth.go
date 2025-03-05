package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/phuslu/log"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/web/webutils"
)

func (m *Middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return session.VerifySession(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.addSessionToCtx(r)
		if err != nil {
			log.Error().Msg(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
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

	reqTenantID := webutils.GetTenantIDFromReq(r)
	if reqTenantID == "" {
		return nil, errors.New(msg.ErrTenantIDNotFoundInRequest)
	}
	tenants, err := webutils.GetUserTenantsFromMetadata(metadata)
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
