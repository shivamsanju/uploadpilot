package handlers

import (
	"net/http"
	"time"

	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
)

var TOKEN_EXPIRY_DURATION = time.Hour * 24 * 30

type apiKeyHandler struct {
	apiKeySvc *services.APIKeyService
}

func NewAPIKeyHandler(apiKeySvc *services.APIKeyService) *apiKeyHandler {
	return &apiKeyHandler{
		apiKeySvc: apiKeySvc,
	}
}

func (h *apiKeyHandler) GetAPIKeys(
	r *http.Request, params interface{}, query interface{}, body interface{},
) ([]models.APIKey, int, error) {
	keys, err := h.apiKeySvc.GetAllAPIKeysForUser(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return keys, http.StatusOK, nil
}

func (h *apiKeyHandler) CreateAPIKey(
	r *http.Request,
	params interface{},
	query interface{},
	body dto.CreateApiKeyData,
) (string, int, error) {
	key, err := h.apiKeySvc.CreateAPIKey(r.Context(), &body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return key, http.StatusOK, nil
}

func (h *apiKeyHandler) RevokeAPIKey(
	r *http.Request,
	params dto.ApiKeyParams,
	query interface{},
	body interface{},
) (bool, int, error) {
	if err := h.apiKeySvc.RevokeAPIKey(r.Context(), params.ApiKeyID); err != nil {
		return false, http.StatusInternalServerError, err
	}
	return true, http.StatusOK, nil
}
