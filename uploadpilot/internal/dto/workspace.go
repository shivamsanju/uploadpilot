package dto

import "github.com/uploadpilot/uploadpilot/internal/db/models"

type WorkspaceUser struct {
	ID    string `json:"userId"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type AddWorkspaceUser struct {
	Email string          `json:"email" validate:"email,required"`
	Role  models.UserRole `json:"role" validate:"required"`
}

type EditUserRole struct {
	Role models.UserRole `json:"role" validate:"required"`
}

type WorkspaceNameID struct {
	Name string `bson:"name" json:"name" validate:"required"`
	ID   string `bson:"_id" json:"id" validate:"required"`
}
