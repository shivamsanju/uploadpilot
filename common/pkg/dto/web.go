package dto

type PaginatedResponse[T any] struct {
	Records      []T   `json:"records"`
	TotalRecords int64 `json:"totalRecords"`
}

type ErrorResponse struct {
	RequestID string `json:"requestID"`
	Message   string `json:"message"`
}
