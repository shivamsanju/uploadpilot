package handlers

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/tenant"
	"github.com/uploadpilot/manager/internal/utils"
)

type tenantHandler struct {
	tenantSvc *tenant.Service
}

func NewTenantHandler(tenantSvc *tenant.Service) *tenantHandler {
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
	if err := utils.ConvertDTOToModel(&body, &tenant); err != nil {
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
