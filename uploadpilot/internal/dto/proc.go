package dto

type EditProcRequest struct {
	Name     string   `json:"name" validate:"required"`
	Triggers []string `json:"triggers" validate:"required"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled" validate:"required"`
}
