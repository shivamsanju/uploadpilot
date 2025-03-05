package models

type Filter struct {
	Field string   `json:"field"`
	Value []string `json:"value"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type PaginationParams struct {
	Offset              int      `json:"offset"`
	Limit               int      `json:"limit"`
	Search              string   `json:"search"`
	CaseSensitiveSearch bool     `json:"caseSensitiveSearch"`
	Filter              []Filter `json:"filter"`
	Sort                Sort     `json:"sort"`
}
