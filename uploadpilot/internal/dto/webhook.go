package dto

type PatchWebhookRequest struct {
	Enabled bool `json:"enabled" validate:"required"`
}
