package dto

type UserContext struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type SessionResponse struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatarUrl"`
	TrialExpiresIn int64  `json:"trialExpiresIn"`
}

type ContextKey string

const (
	UserIDContextKey ContextKey = "id"
	EmailContextKey  ContextKey = "email"
	NameContextKey   ContextKey = "name"
)
