package services

import (
	"context"
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/web/webutils"
)

type TenantService struct {
	acm        *rbac.AccessManager
	tenantRepo *repo.TenantRepo
}

func NewTenantService(accessManager *rbac.AccessManager, tenantRepo *repo.TenantRepo) *TenantService {
	return &TenantService{
		acm:        accessManager,
		tenantRepo: tenantRepo,
	}
}

func (s *TenantService) GetTenantDetails(ctx context.Context, tenantID string) (*models.Tenant, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.acm.CheckAccess(session.Sub, tenantID, "", rbac.Reader) {
		return nil, err
	}

	return s.tenantRepo.Get(ctx, tenantID)
}

func (s *TenantService) OnboardTenant(ctx context.Context, tenant *models.Tenant) error {
	sessionContainer := session.GetSessionFromRequestContext(ctx)
	userID := sessionContainer.GetUserID()
	if userID == "" {
		return errors.New(msg.ErrUserInfoNotFoundInRequest)
	}

	tenant.OwnerID = userID
	tenant.Status = models.TenantStatusActive

	if err := s.tenantRepo.Create(ctx, tenant, s.tenantOnboardingCallback); err != nil {
		return err
	}

	return nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.acm.CheckAccess(session.Sub, tenantID, "", rbac.Admin) {
		return err
	}

	return s.tenantRepo.Delete(ctx, tenantID)
}

func (s *TenantService) tenantOnboardingCallback(tenant *models.Tenant) error {
	met, err := usermetadata.GetUserMetadata(tenant.OwnerID)
	if err != nil {
		return err
	}

	tenants, err := webutils.GetUserTenantsFromMetadata(met)
	if err != nil {
		return err
	}

	if err := s.acm.AddAccess(tenant.OwnerID, tenant.ID, "", rbac.Admin); err != nil {
		return err
	}

	tenants[tenant.ID] = dto.TenantMetadata{Name: tenant.Name, ID: tenant.ID}

	_, err = usermetadata.UpdateUserMetadata(tenant.OwnerID, map[string]interface{}{
		dto.UserMetadataTenantKey: tenants,
		dto.ActiveTenantIDKey:     tenant.ID,
	})
	return err
}
