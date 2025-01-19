package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/uploadpilot/uploadpilot/internal/auth"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"github.com/uploadpilot/uploadpilot/web/dto"
	"golang.org/x/net/context"
)

type authHandler struct {
	userRepo db.UserRepo
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
		token, err := auth.GetSignedToken(w, &user)
		if err != nil {
			redrectWithError(w, r, err)
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
		redrectWithError(w, r, err)
		return
	}
	existingProvider, err := h.userRepo.CheckUserProvider(r.Context(), user.Email)
	if err != nil {
		redrectWithError(w, r, err)

		return
	}
	if existingProvider != "" && existingProvider != provider {
		err := fmt.Errorf("user with email %s already exists with provider %s. please login with the same provider you used earlier", user.Email, existingProvider)
		redrectWithError(w, r, err)
		return
	}
	if existingProvider == "" {
		h.userRepo.CreateUser(r.Context(), mapUser(&user))
	}
	token, err := auth.GetSignedToken(w, &user)
	if err != nil {
		redrectWithError(w, r, err)
		return
	}
	redirectUri := getTokenRedirectURI(token)
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	auth.RemoveBearerTokenInCookie(w)
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
	userID := r.Header.Get("userId")
	cb, err := h.userRepo.GetUserByID(r.Context(), userID)
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
	infra.Log.Infof("mapping user: %+v", user)
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

func redrectWithError(w http.ResponseWriter, r *http.Request, err error) {
	redirectUri := config.FrontendURI + "/error?error=" + err.Error()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}
