package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	webmodels "github.com/shivamsanju/uploader/internal/web/models"
	"github.com/shivamsanju/uploader/internal/web/utils"
	"github.com/shivamsanju/uploader/pkg/globals"
)

type userHandler struct {
	userRepo repo.UserRepo
}

func NewUserHandler() *userHandler {
	return &userHandler{
		userRepo: repo.NewUserRepo(),
	}
}

var validate = validator.New()

func (h *userHandler) Signup(w http.ResponseWriter, r *http.Request) {
	signupReq := &webmodels.SignupRequest{}
	if err := render.DecodeJSON(r.Body, signupReq); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(signupReq); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}
	if signupReq.RootPassword != globals.RootPassword {
		utils.HandleHttpError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid root password"))
		return
	}
	if signupReq.Password != signupReq.ConfirmPassword {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("password and confirm password does not match"))
		return
	}
	userExists, _ := h.userRepo.CheckUserExists(r.Context(), signupReq.Email)
	if userExists {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}
	user := &models.User{
		FirstName: signupReq.FirstName,
		LastName:  signupReq.LastName,
		Email:     signupReq.Email,
	}
	user.Password = utils.HashPassword(signupReq.Password)
	token, refreshToken, err := utils.GenerateTokens(user.Email, user.FirstName, user.LastName, user.ID.Hex())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	user.Token = token
	user.RefreshToken = refreshToken
	_, err = h.userRepo.CreateUser(r.Context(), user)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, user)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginReq := webmodels.LoginRequest{}
	if err := render.DecodeJSON(r.Body, &loginReq); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(loginReq); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	user, err := h.userRepo.GetUserByEmail(r.Context(), loginReq.Email)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("user with email %s not found", loginReq.Email))
		return
	}

	if !utils.VerifyPassword(loginReq.Password, user.Password) {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}

	token, refreshToken, err := utils.GenerateTokens(user.Email, user.FirstName, user.LastName, user.ID.Hex())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	user.Token = token
	user.RefreshToken = refreshToken
	render.JSON(w, r, user)
}

func (h *userHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	cb, err := h.userRepo.GetUserByEmail(r.Context(), email)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}
