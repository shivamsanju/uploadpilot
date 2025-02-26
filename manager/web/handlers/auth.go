package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/phuslu/log"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/auth"
	"github.com/uploadpilot/manager/internal/svc/user"
	"github.com/uploadpilot/manager/internal/utils"
	"golang.org/x/net/context"
)

var TOKEN_EXPIRY_DURATION = time.Hour * 24 * 30

type authHandler struct {
	userSvc *user.Service
	authSvc *auth.Service
}

func NewAuthHandler(userSvc *user.Service, authSvc *auth.Service) *authHandler {
	return &authHandler{
		userSvc: userSvc,
		authSvc: authSvc,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		usr, err := h.userSvc.GetUserByEmail(r.Context(), user.Email)
		if err != nil {
			redirectWithError(w, r, fmt.Errorf("failed to get user"))
			return
		}
		token, err := h.authSvc.GenerateJWTToken(w, usr, TOKEN_EXPIRY_DURATION)
		if err != nil {
			redirectWithError(w, r, err)
			return
		}
		redirectUri := getTokenRedirectURI(token)
		http.Redirect(w, r, redirectUri, http.StatusSeeOther)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *authHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	log.Info().Msgf("handling callback for provider: %s", provider)
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		redirectWithError(w, r, err)
		return
	}

	usr, err := h.userSvc.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		redirectWithError(w, r, err)
		return
	}

	if usr == nil || usr.ID == "" {
		usr = mapUser(&user)
		if err := h.userSvc.CreateUser(r.Context(), usr); err != nil {
			redirectWithError(w, r, err)
			return
		}
	}

	token, err := h.authSvc.GenerateJWTToken(w, usr, TOKEN_EXPIRY_DURATION)
	if err != nil {
		redirectWithError(w, r, err)
		return
	}
	redirectUri := getTokenRedirectURI(token)
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	redirectUri := getLogoutRedirectURI()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *authHandler) LogoutProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.Logout(w, r)
	redirectUri := getLogoutRedirectURI()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *authHandler) GetAPIKeys(
	r *http.Request, params interface{}, query interface{}, body interface{},
) ([]models.APIKey, int, error) {
	keys, err := h.authSvc.GetAllApiKeysForUser(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return keys, http.StatusOK, nil
}

func (h *authHandler) CreateAPIKey(
	r *http.Request,
	params interface{},
	query interface{},
	body dto.CreateApiKeyData,
) (string, int, error) {
	log.Debug().Msgf("creating api key with data %+v", body)
	key, err := h.authSvc.CreateApiKey(r.Context(), &body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return key, http.StatusOK, nil
}

func (h *authHandler) RevokeAPIKey(
	r *http.Request,
	params dto.ApiKeyParams,
	query interface{},
	body interface{},
) (bool, int, error) {
	if err := h.authSvc.RevokeApiKey(r.Context(), params.ApiKeyID); err != nil {
		return false, http.StatusInternalServerError, err
	}
	return true, http.StatusOK, nil
}

func (h *authHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(dto.UserIDContextKey).(string)
	cb, err := h.userSvc.GetUserDetails(r.Context(), userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	trialsExpiresIn := time.Until(cb.TrialEndsAt).Hours()
	render.JSON(w, r, &dto.SessionResponse{
		Name:           *cb.Name,
		Email:          cb.Email,
		AvatarURL:      *cb.AvatarURL,
		TrialExpiresIn: int64(trialsExpiresIn),
	})
}

func getTokenRedirectURI(token string) string {
	return config.FrontendURI + "/auth?uploadpilottoken=" + token
}

func getLogoutRedirectURI() string {
	return config.FrontendURI + "/auth"
}

func redirectWithError(w http.ResponseWriter, r *http.Request, err error) {
	redirectUri := config.FrontendURI + "/error?error=" + err.Error()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func mapUser(user *goth.User) *models.User {
	return &models.User{
		Email:         user.Email,
		Provider:      &user.Provider,
		Name:          &user.Name,
		FirstName:     &user.FirstName,
		LastName:      &user.LastName,
		NickName:      &user.NickName,
		AvatarURL:     &user.AvatarURL,
		Location:      &user.Location,
		Description:   &user.Description,
		TrialStartsAt: time.Now(),
		TrialEndsAt:   time.Now().Add(time.Hour * 14 * 24),
	}
}
