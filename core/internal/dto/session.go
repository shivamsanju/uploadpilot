package dto

type SessionCtxKeyType string

const SessionCtxKey SessionCtxKeyType = "user_info"

type Session struct {
	UserID   string                 `json:"id"`
	Email    string                 `json:"email"`
	Metadata map[string]interface{} `json:"metadata"`
}
