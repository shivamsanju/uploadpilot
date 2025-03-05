package dto

type TenantOnboardingRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=25,alphanumspace"`
	ContactEmail string `json:"contactEmail" validate:"required,email"`
	Phone        string `json:"phone" validate:"max=20"`
	Address      string `json:"address" validate:"max=500"`
	Industry     string `json:"industry" validate:"max=50,alphanumspace"`
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
