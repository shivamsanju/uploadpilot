package models

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignupRequest struct {
	FirstName       string `json:"firstName" validate:"required,min=2,max=100"`
	LastName        string `json:"lastName" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"email,required"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6"`
}
