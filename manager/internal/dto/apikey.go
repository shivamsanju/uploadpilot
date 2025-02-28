package dto

import (
	"time"

	"github.com/uploadpilot/go-core/db/pkg/dtypes"
	"github.com/uploadpilot/go-core/db/pkg/models"
)

type UserClaims struct {
	UserID string
	Email  string
}

type APIKeyPerm struct {
	Scope     string
	ResouceID string
	Perm      models.APIKeyPermissionType
}

type CreateApiKeyData struct {
	Name        string                    `json:"name" validate:"required,max=20"`
	ExpiresAt   time.Time                 `json:"expiresAt" validate:"required,future"`
	Scopes      dtypes.StringArray        `json:"scopes"`
	Permissions []models.APIKeyPermission `json:"permissions"`
}

type ApiKeyLimitedDetails struct {
	ID          string    `json:"id"`
	Comment     string    `json:"comment"`
	WorkspaceId string    `json:"workspaceId"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Revoked     bool      `json:"revoked"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedBy   string    `json:"updatedBy"`
}
