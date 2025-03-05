package dto

var UserAttributesKey = "userAttributes"

type UserAttributes struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Theme  string `json:"theme,omitempty"`
}

type UserDetailsResponse struct {
	Email        string            `json:"email"`
	Name         string            `json:"name,omitempty"`
	Avatar       *string           `json:"avatar,omitempty"`
	Theme        *string           `json:"theme,omitempty"`
	Tenants      map[string]string `json:"tenants,omitempty"`
	ActiveTenant *string           `json:"activeTenant,omitempty"`
}
