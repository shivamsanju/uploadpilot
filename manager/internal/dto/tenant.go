package dto

type TenantOnboardingRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=25,alphanum"`
	ContactEmail string `json:"contactEmail" validate:"required,email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Industry     string `json:"industry"`
	CompanyName  string `json:"companyName"`
	Role         string `json:"role" validate:"required"`
}

var UserMetadataTenantKey = "tenants"
var ActiveTenantIDKey = "activeTenantID"

type TenantMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SetActiveTenant struct {
	TenantID string `json:"tenantId" validate:"required,uuid"`
}
