package dto

type PaginatedResponse[T any] struct {
	Records      []T   `json:"records"`
	TotalRecords int64 `json:"totalRecords"`
}

type SessionResponse struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}
