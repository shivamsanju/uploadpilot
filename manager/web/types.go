package web

import "github.com/uploadpilot/manager/internal/dto"

type APIAuthType string

const (
	APIAuthTypeBearer APIAuthType = "Bearer"
	APIAuthTypeAPIKey APIAuthType = "APIKey"
)

type APIPermission struct {
	AllowedAuthTypes []APIAuthType
	Permissions      []dto.APIKeyPerm
}
