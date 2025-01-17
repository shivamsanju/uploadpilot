package models

type PaginatedResponse[T any] struct {
	Records      []T   `json:"records"`
	TotalRecords int64 `json:"totalRecords"`
}
