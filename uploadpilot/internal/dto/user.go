package dto

type ApiUser struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type SessionResponse struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatarUrl"`
	TrialExpiresIn int64  `json:"trialExpiresIn"`
}
