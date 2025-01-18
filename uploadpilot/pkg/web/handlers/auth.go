package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/uploadpilot/uploadpilot/pkg/auth"
	"github.com/uploadpilot/uploadpilot/pkg/db"
	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	"github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/utils"
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
			utils.HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}
		globals.Log.Infof("user logged in: %+v", user)
		redirectUri := getTokenRedirectURI(token)
		http.Redirect(w, r, redirectUri, http.StatusSeeOther)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *authHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	globals.Log.Infof("handling callback for provider: %s", provider)
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		utils.HandleHttpError(w, r, http.StatusUnauthorized, err)
		return
	}
	exists, err := h.userRepo.CheckUserExists(r.Context(), user.UserID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if !exists {
		h.userRepo.CreateUser(r.Context(), mapUser(&user))
	}
	token, err := auth.GetSignedToken(w, &user)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	globals.Log.Infof("user logged in: %+v", user)
	redirectUri := getTokenRedirectURI(token)
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	auth.RemoveBearerTokenInCookie(w)
	http.Redirect(w, r, globals.FrontendURI+"/auth", http.StatusSeeOther)
}

func (h *authHandler) LogoutProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *authHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userId")
	cb, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func mapUser(user *goth.User) *models.User {
	globals.Log.Infof("mapping user: %+v", user)
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
	return globals.FrontendURI + "/auth?uploadpilottoken=" + token
}
