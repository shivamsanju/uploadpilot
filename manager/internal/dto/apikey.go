package dto

import (
	"time"

	"github.com/uploadpilot/manager/internal/db/models"
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

type WorkspaceApiPerm struct {
	ID     string `json:"id"`
	Read   bool   `json:"read"`
	Manage bool   `json:"manage"`
	Upload bool   `json:"upload"`
}
type CreateApiKeyData struct {
	Name           string             `json:"name" validate:"required,min=3,max=25,alphanumspace"`
	ExpiresAt      time.Time          `json:"expiresAt" validate:"required,future"`
	TenantRead     bool               `json:"tenantRead"`
	TenantManage   bool               `json:"tenantManage"`
	WorkspacePerms []WorkspaceApiPerm `json:"workspacePerms"`
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
