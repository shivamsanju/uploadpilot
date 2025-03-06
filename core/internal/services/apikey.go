package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	"github.com/uploadpilot/core/config"
	"github.com/uploadpilot/core/internal/db/errs"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/pkg/utils"
	"github.com/uploadpilot/core/pkg/vault"
	"github.com/uploadpilot/core/web/webutils"
)

type APIKeyService struct {
	apiKeyRepo *repo.APIKeyRepo
	kms        vault.KMS
	apiKeySalt string
}

func NewAPIKeyService(apiKeyRepo *repo.APIKeyRepo, kms vault.KMS) *APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
		kms:        kms,
		apiKeySalt: config.AppConfig.ApiKeyEncryptionKey,
	}
}

func (s *APIKeyService) GetAllAPIKeysForUser(ctx context.Context, tenantID string) ([]models.APIKey, error) {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.apiKeyRepo.GetAllApiKeys(ctx, user.UserID, tenantID)
}

func (s *APIKeyService) GetAPIKeyInfo(ctx context.Context, id string) (*models.APIKey, error) {
	return s.apiKeyRepo.GetApiKeyDetailsByID(ctx, id)
}

func (s *APIKeyService) CreateAPIKey(ctx context.Context, tenantID string, data *dto.CreateApiKeyData) (string, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return "", err
	}

	perms, scope, err := s.getScopeAndPerm(tenantID, session, data)
	if err != nil {
		return "", err
	}

	newKey := "up-" + utils.GenerateRandomAlphaNumericString(64) + data.ExpiresAt.Format("20060102150405")

	hashedKey, err := s.kms.Hash(newKey, s.apiKeySalt)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash api key")
		return "", errors.New(msg.ErrApiKeyCreateFailed)
	}

	apiKey := &models.APIKey{
		Name:        data.Name,
		UserID:      session.UserID,
		TenantID:    tenantID,
		ApiKeyHash:  hashedKey,
		ExpiresAt:   &data.ExpiresAt,
		Scopes:      scope,
		Permissions: perms,
	}

	if err := s.apiKeyRepo.CreateApiKey(ctx, apiKey); err != nil {
		return "", err
	}

	return newKey, nil
}

func (s *APIKeyService) RevokeAPIKey(ctx context.Context, id string) error {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	apiKey, err := s.apiKeyRepo.GetApiKeyDetailsByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	apiKey.RevokedAt = &now
	apiKey.RevokedBy = &user.UserID

	return s.apiKeyRepo.Update(ctx, apiKey)
}

func (s *APIKeyService) ValidateAPIKey(ctx context.Context, apiKey string, perms ...dto.APIKeyPerm) (*models.APIKey, error) {
	if ok := s.isValidAPIKeyFormat(apiKey); !ok {
		return nil, errors.New(msg.ErrInvalidAPIKey)
	}

	apiKeyHash, err := s.kms.Hash(apiKey, s.apiKeySalt)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash api key")
		return nil, errors.New(msg.ErrUnexpected)
	}

	apiKeyDetails, err := s.apiKeyRepo.GetApiKeyDetails(ctx, apiKeyHash)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, errors.New(msg.ErrInvalidAPIKey)
		}
		return nil, errors.New(msg.ErrUnexpected)
	}

	if apiKeyDetails.RevokedAt != nil && apiKeyDetails.RevokedAt.Before(time.Now()) {
		log.Debug().Str("api_key", apiKey).Msg("api key revoked")
		return nil, errors.New(msg.ErrRevokedAPIKey)
	}

	if apiKeyDetails.ExpiresAt != nil && apiKeyDetails.ExpiresAt.Before(time.Now()) {
		log.Debug().Str("api_key", apiKey).Msg("api key expired")
		return nil, errors.New(msg.ErrExpiredAPIKey)
	}

	hasPerm := s.verifyAPIKeyPermissions(apiKeyDetails, perms...)

	if !hasPerm {
		return nil, errors.New(msg.ErrInvalidAPIKey)
	}

	return apiKeyDetails, nil
}

func (s *APIKeyService) isValidAPIKeyFormat(apiKey string) bool {
	if apiKey == "" || !strings.HasPrefix(apiKey, "up-") {
		return false
	}

	if len(apiKey) < 14 || len(apiKey) > 74 {
		return false
	}

	timePart := apiKey[len(apiKey)-14:]
	tp, err := strconv.ParseInt(timePart, 10, 64)
	if err != nil || time.Now().Unix() > tp {
		return false
	}
	return true
}

func (s *APIKeyService) verifyAPIKeyPermissions(apiKey *models.APIKey, perms ...dto.APIKeyPerm) bool {
	for _, perm := range perms {
		if !slices.Contains(apiKey.Scopes, perm.Scope) {
			return false
		}
		for _, apiPerm := range apiKey.Permissions {
			log.Debug().Str("Got", fmt.Sprintf("%s:%s", apiPerm.ResourceID, apiPerm.Permission)).Str("Expected", fmt.Sprintf("%s:%s", perm.ResouceID, perm.Perm)).Msg("verifying api key permissions")

			if perm.ResouceID == apiPerm.ResourceID && perm.Perm == apiPerm.Permission {
				return true
			}
		}
	}
	return false
}

func (s *APIKeyService) getScopeAndPerm(tenantID string, session *dto.Session, data *dto.CreateApiKeyData) ([]models.APIKeyPermission, []string, error) {
	scopes := make(map[string]struct{})
	var perms []models.APIKeyPermission

	if data.TenantRead {
		scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeTenant, models.APIKeyPermissionTypeRead)] = struct{}{}
		perms = append(perms, models.APIKeyPermission{
			ResourceID:   tenantID,
			ResourceType: models.APIPermResourceTypeTenant,
			Permission:   models.APIKeyPermissionTypeRead,
		})
	}
	if data.TenantManage {
		scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeTenant, models.APIKeyPermissionTypeManage)] = struct{}{}
		perms = append(perms, models.APIKeyPermission{
			ResourceID:   tenantID,
			ResourceType: models.APIPermResourceTypeTenant,
			Permission:   models.APIKeyPermissionTypeManage,
		})
	}

	for _, perm := range data.WorkspacePerms {
		if perm.Read {
			scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeWorkspace, models.APIKeyPermissionTypeRead)] = struct{}{}
			perms = append(perms, models.APIKeyPermission{
				ResourceID:   perm.ID,
				ResourceType: models.APIPermResourceTypeWorkspace,
				Permission:   models.APIKeyPermissionTypeRead,
			})
		}
		if perm.Manage {
			scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeWorkspace, models.APIKeyPermissionTypeManage)] = struct{}{}
			perms = append(perms, models.APIKeyPermission{
				ResourceID:   perm.ID,
				ResourceType: models.APIPermResourceTypeWorkspace,
				Permission:   models.APIKeyPermissionTypeManage,
			})
		}
		if perm.Upload {
			scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeWorkspace, models.APIKeyPermissionTypeUpload)] = struct{}{}
			perms = append(perms, models.APIKeyPermission{
				ResourceID:   perm.ID,
				ResourceType: models.APIPermResourceTypeWorkspace,
				Permission:   models.APIKeyPermissionTypeUpload,
			})
		}
	}

	if len(scopes) == 0 {
		return nil, nil, errors.New(msg.ErrNoScopeInAPIKeyCreateRequest)
	}

	scopesSlice := make([]string, 0, len(scopes))
	for k := range scopes {
		scopesSlice = append(scopesSlice, k)
	}

	return perms, scopesSlice, nil
}
