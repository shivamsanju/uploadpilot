package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type AuthHandler interface {
	GetUserInfo(w http.ResponseWriter, r *http.Request)
}

type authHandler struct{}

func NewAuthHandler() AuthHandler {
	return &authHandler{}
}

func (h *authHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	userId := sessionContainer.GetUserID()
	sessionData, err := thirdparty.GetUserByID(userId)
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	metadata, err := usermetadata.GetUserMetadata(userId)
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	render.JSON(w, r, map[string]interface{}{
		"email": sessionData.Email,
		"image": metadata["picture"],
	})
}
