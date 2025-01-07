package models

type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}
