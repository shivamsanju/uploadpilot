package dto

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/uploadpilot/go-core/db/pkg/models"
)

type UserClaims struct {
	UserID string
	Email  string
	Name   string
}

type JWTClaims struct {
	*UserClaims
	*jwt.StandardClaims
}

type APIKeyClaims struct {
	*UserClaims
	CanReadAcc   bool
	CanManageAcc bool
	Permissions  []models.APIKeyPerm
}

type CreateApiKeyData struct {
	Name         string              `json:"name" validate:"required,max=20"`
	ExpiresAt    time.Time           `json:"expiresAt" validate:"required,future"`
	CanManageAcc bool                `json:"canManageAcc"`
	CanReadAcc   bool                `json:"canReadAcc"`
	Permissions  []models.APIKeyPerm `json:"permissions"`
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
