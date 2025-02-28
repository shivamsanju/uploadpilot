package web

import "github.com/uploadpilot/manager/internal/dto"

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
