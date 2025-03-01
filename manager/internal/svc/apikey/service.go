package apikey

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/db/errs"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/db/repo"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/manager/internal/utils"
	"github.com/uploadpilot/manager/internal/vault"
)

type Service struct {
	apiKeyRepo *repo.APIKeyRepo
	kms        vault.KMS
	apiKeySalt string
}

func NewService(apiKeyRepo *repo.APIKeyRepo, kms vault.KMS) *Service {
	return &Service{
		apiKeyRepo: apiKeyRepo,
		kms:        kms,
		apiKeySalt: config.ApiKeyEncryptionSalt,
	}
}

func (s *Service) GetAllAPIKeysForUser(ctx context.Context) ([]models.APIKey, error) {
	user, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.apiKeyRepo.GetAllApiKeys(ctx, user.UserID)
}

func (s *Service) GetAPIKeyInfo(ctx context.Context, id string) (*models.APIKey, error) {
	return s.apiKeyRepo.GetApiKeyDetailsByID(ctx, id)
}

func (s *Service) CreateAPIKey(ctx context.Context, data *dto.CreateApiKeyData) (string, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return "", err
	}

	perms, scope, err := s.getScopeAndPerm(session, data)
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
		TenantID:    session.TenantID,
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

func (s *Service) RevokeAPIKey(ctx context.Context, id string) error {
	user, err := utils.GetSessionFromCtx(ctx)
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

func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string, perms ...dto.APIKeyPerm) (*models.APIKey, error) {
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

func (s *Service) isValidAPIKeyFormat(apiKey string) bool {
	if apiKey == "" && !strings.HasPrefix(apiKey, "up-") {
		return false
	}

	timePart := apiKey[len(apiKey)-14:]
	tp, err := strconv.ParseInt(timePart, 10, 64)
	if err != nil || time.Now().Unix() > tp {
		return false
	}
	return true
}

func (s *Service) verifyAPIKeyPermissions(apiKey *models.APIKey, perms ...dto.APIKeyPerm) bool {
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

func (s *Service) getScopeAndPerm(session *dto.Session, data *dto.CreateApiKeyData) ([]models.APIKeyPermission, []string, error) {
	scopes := make(map[string]struct{})
	var perms []models.APIKeyPermission

	if data.TenantRead {
		scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeTenant, models.APIKeyPermissionTypeRead)] = struct{}{}
		perms = append(perms, models.APIKeyPermission{
			ResourceID:   session.TenantID,
			ResourceType: models.APIPermResourceTypeTenant,
			Permission:   models.APIKeyPermissionTypeRead,
		})
	}
	if data.TenantManage {
		scopes[fmt.Sprintf("%s:%s", models.APIPermResourceTypeTenant, models.APIKeyPermissionTypeManage)] = struct{}{}
		perms = append(perms, models.APIKeyPermission{
			ResourceID:   session.TenantID,
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
