package web

import (
	"github.com/uploadpilot/manager/internal/dto"
)

type APIAuthType string

const (
	APIAuthTypeBearer APIAuthType = "Bearer"
	APIAuthTypeAPIKey APIAuthType = "APIKey"
)

type APIPermission struct {
	AllowedAuthTypes []APIAuthType
	Permissions      []dto.APIKeyPerm
}

var BearerTenantReadAccess = APIPermission{
	AllowedAuthTypes: []APIAuthType{APIAuthTypeBearer},
	Permissions: []dto.APIKeyPerm{
		{
			Scope:     "tenant",
			ResouceID: "tenant-id",
			Perm:      "read",
		},
	},
}
