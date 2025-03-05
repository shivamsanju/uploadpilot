package handlers

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
)

type tenantHandler struct {
	tenantSvc *services.TenantService
}

func NewTenantHandler(tenantSvc *services.TenantService) *tenantHandler {
	return &tenantHandler{
		tenantSvc: tenantSvc,
	}
}

func (h *tenantHandler) OnboardTenant(
	r *http.Request,
	params interface{},
	query interface{},
	body dto.TenantOnboardingRequest,
) (*string, int, error) {
	var tenant models.Tenant
	if err := copier.Copy(&tenant, &body); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}
	if err := h.tenantSvc.OnboardTenant(r.Context(), &tenant); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &tenant.ID, http.StatusOK, nil
}

func (h *tenantHandler) SetActiveTenant(
	r *http.Request,
	params interface{},
	query interface{},
	body dto.SetActiveTenant,
) (string, int, error) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userID := sessionContainer.GetUserID()

	if _, err := usermetadata.UpdateUserMetadata(userID, map[string]interface{}{
		dto.ActiveTenantIDKey: body.TenantID,
	}); err != nil {
		return "", http.StatusBadRequest, err
	}

	return body.TenantID, http.StatusOK, nil
}
