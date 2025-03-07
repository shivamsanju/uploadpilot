package services

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/pkg/utils"
	"github.com/uploadpilot/core/pkg/vault"
	"github.com/uploadpilot/core/web/webutils"
)

type APIKeyService struct {
	accessManager *rbac.AccessManager
	apiKeyRepo    *repo.APIKeyRepo
	kms           vault.KMS
	apiKeySalt    string
}

func NewAPIKeyService(accessManager *rbac.AccessManager, apiKeyRepo *repo.APIKeyRepo, kms vault.KMS) *APIKeyService {
	return &APIKeyService{
		accessManager: accessManager,
		apiKeyRepo:    apiKeyRepo,
		kms:           kms,
		apiKeySalt:    config.AppConfig.ApiKeyEncryptionKey,
	}
}

func (s *APIKeyService) GetAllAPIKeysInTenant(ctx context.Context, tenantID string) ([]models.APIKey, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, "", rbac.Admin) {
		return nil, fmt.Errorf(msg.ErrAccessDenied)
	}

	return s.apiKeyRepo.GetAllApiKeysInTenant(ctx, tenantID)
}

func (s *APIKeyService) GetAPIKeyInfo(ctx context.Context, tenantID, id string) (*models.APIKey, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, "", rbac.Admin) {
		return nil, fmt.Errorf(msg.ErrAccessDenied)
	}

	return s.apiKeyRepo.GetApiKeyDetailsByID(ctx, id)
}

func (s *APIKeyService) CreateAPIKey(ctx context.Context, tenantID string, data *dto.CreateApiKeyData) (string, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return "", err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, "", rbac.Admin) {
		return "", fmt.Errorf(msg.ErrAccessDenied)
	}

	newKey := "up-" + utils.GenerateRandomAlphaNumericString(64) + data.ExpiresAt.Format("20060102150405")
	hashedKey, err := s.kms.Hash(newKey, s.apiKeySalt)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash api key")
		return "", errors.New(msg.ErrApiKeyCreateFailed)
	}

	apiKey := &models.APIKey{
		Name:       data.Name,
		UserID:     session.UserID,
		TenantID:   tenantID,
		ApiKeyHash: hashedKey,
		ExpiresAt:  &data.ExpiresAt,
	}

	if err := s.apiKeyRepo.CreateApiKey(ctx, apiKey, func() error {
		return s.addAccessToAPIKey(hashedKey, tenantID, data)
	}); err != nil {
		return "", err
	}

	return newKey, nil
}

func (s *APIKeyService) RevokeAPIKey(ctx context.Context, tenantID, id string) error {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.accessManager.CheckAccess(user.UserID, tenantID, "", rbac.Admin) {
		return fmt.Errorf(msg.ErrAccessDenied)
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

func (s *APIKeyService) VerifyAPIKey(ctx context.Context, apiKey string) (*models.APIKey, error) {
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

	return apiKeyDetails, nil
}

func (s *APIKeyService) addAccessToAPIKey(apiKeyHash, tenantID string, data *dto.CreateApiKeyData) error {
	if data.TenantRead {
		if err := s.accessManager.AddAccess(apiKeyHash, tenantID, "", rbac.Reader); err != nil {
			return err
		}
	}

	for _, perm := range data.WorkspacePerms {
		if perm.Manage {
			if err := s.accessManager.AddAccess(apiKeyHash, tenantID, perm.ID, rbac.Admin); err != nil {
				return err
			}
		}
		if perm.Read {
			if err := s.accessManager.AddAccess(apiKeyHash, tenantID, perm.ID, rbac.Reader); err != nil {
				return err
			}
		}
		if perm.Upload {
			if err := s.accessManager.AddAccess(apiKeyHash, tenantID, perm.ID, rbac.Uploader); err != nil {
				return err
			}
		}
	}

	return nil
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
