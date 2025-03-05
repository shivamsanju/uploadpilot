package webutils

import (
	"fmt"

	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
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

var TenantReadAccess = APIPermission{
	AllowedAuthTypes: []APIAuthType{APIAuthTypeBearer},
	Permissions: []dto.APIKeyPerm{
		{
			Scope:     fmt.Sprintf("%s:%s", models.APIPermResourceTypeTenant, models.APIKeyPermissionTypeRead),
			ResouceID: "<tenantId>",
			Perm:      "read",
		},
	},
}

var WorkspaceUploadAccess = APIPermission{
	AllowedAuthTypes: []APIAuthType{APIAuthTypeAPIKey, APIAuthTypeBearer},
	Permissions: []dto.APIKeyPerm{
		{
			Scope:     fmt.Sprintf("%s:%s", models.APIPermResourceTypeWorkspace, models.APIKeyPermissionTypeUpload),
			ResouceID: "<workspaceId>",
			Perm:      "upload",
		},
	},
}
