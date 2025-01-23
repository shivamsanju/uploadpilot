package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/uploadpilot/uploadpilot/internal/auth"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"golang.org/x/net/context"
)

var TOKEN_EXPIRY_DURATION = time.Hour * 24 * 30

type authHandler struct {
	userRepo *db.UserRepo
}

func NewAuthHandler() *authHandler {
	return &authHandler{
		userRepo: db.NewUserRepo(),
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		token, err := auth.GenerateToken(w, &user, TOKEN_EXPIRY_DURATION)
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
	infra.Log.Infof("handling callback for provider: %s", provider)
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		redirectWithError(w, r, err)
		return
	}
	existingProvider, err := h.userRepo.GetProvider(r.Context(), user.Email)
	if err != nil {
		redirectWithError(w, r, err)

		return
	}
	if existingProvider != "" && existingProvider != provider {
		err := fmt.Errorf("user with email %s already exists with provider %s. please login with the same provider you used earlier", user.Email, existingProvider)
		redirectWithError(w, r, err)
		return
	}
	if existingProvider == "" {
		h.userRepo.Create(r.Context(), mapUser(&user))
	}
	token, err := auth.GenerateToken(w, &user, TOKEN_EXPIRY_DURATION)
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

func (h *authHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(string)
	cb, err := h.userRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, &dto.SessionResponse{
		Name:      cb.Name,
		Email:     cb.Email,
		AvatarURL: cb.AvatarURL,
	})
}

func mapUser(user *goth.User) *models.User {
	return &models.User{
		UserID:        user.UserID,
		Email:         user.Email,
		Provider:      user.Provider,
		Name:          user.Name,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		NickName:      user.NickName,
		AvatarURL:     user.AvatarURL,
		Location:      user.Location,
		Description:   user.Description,
		IsUserBanned:  false,
		EmailVerified: true,
	}
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
