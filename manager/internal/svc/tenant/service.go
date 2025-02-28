package tenant

import (
	"context"
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/manager/internal/utils"
)

type Service struct {
	tenantRepo *repo.TenantRepo
}

func NewService(tenantRepo *repo.TenantRepo) *Service {
	return &Service{tenantRepo: tenantRepo}
}

func (s *Service) GetTenantDetails(ctx context.Context, tenantID string) (*models.Tenant, error) {
	return s.tenantRepo.Get(ctx, tenantID)
}

func (s *Service) OnboardTenant(ctx context.Context, tenant *models.Tenant) error {
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

func (s *Service) DeleteTenant(ctx context.Context, tenantID string) error {
	return s.tenantRepo.Delete(ctx, tenantID)
}

func (s *Service) tenantOnboardingCallback(tenant *models.Tenant) error {
	met, err := usermetadata.GetUserMetadata(tenant.OwnerID)
	if err != nil {
		return err
	}

	tenants, err := utils.GetUserTenantsFromMetadata(met)
	if err != nil {
		return err
	}

	tenants[tenant.ID] = dto.TenantMetadata{Name: tenant.Name, ID: tenant.ID}

	_, err = usermetadata.UpdateUserMetadata(tenant.OwnerID, map[string]interface{}{
		dto.UserMetadataTenantKey: tenants,
		dto.ActiveTenantIDKey:     tenant.ID,
	})
	return err
}
