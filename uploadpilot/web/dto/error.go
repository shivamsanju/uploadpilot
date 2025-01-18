package dto

type ErrorResponse struct {
	RequestID string `json:"requestID"`
	Message   string `json:"message"`
}
