package dto

import "github.com/uploadpilot/uploadpilot/internal/db/models"

type AddUserToWorkspaceRequest struct {
	Email string          `json:"email" validate:"email,required"`
	Role  models.UserRole `json:"role" validate:"required"`
}

type EditUserInWorkspaceRequest struct {
	Role models.UserRole `json:"role" validate:"required"`
}
