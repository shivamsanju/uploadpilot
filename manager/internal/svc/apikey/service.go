package apikey

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	commonutils "github.com/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/go-core/common/vault"
	"github.com/uploadpilot/go-core/db/pkg/errs"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/manager/internal/utils"
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

func (s *Service) GetAPIKeyInfo(ctx context.Context, key string) (*models.APIKey, error) {
	return s.apiKeyRepo.GetApiKeyDetails(ctx, key)
}

func (s *Service) CreateAPIKey(ctx context.Context, data *dto.CreateApiKeyData) (string, error) {
	user, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return "", err
	}
	newKey := "up-" + commonutils.GenerateRandomAlphaNumericString(64) + data.ExpiresAt.Format("20060102150405")

	hashedKey, err := s.kms.Hash(newKey, s.apiKeySalt)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash api key")
		return "", errors.New(msg.ErrApiKeyCreateFailed)
	}

	apiKey := &models.APIKey{
		Name:        data.Name,
		UserID:      user.UserID,
		ApiKeyHash:  hashedKey,
		ExpiresAt:   &data.ExpiresAt,
		Scopes:      data.Scopes,
		Permissions: data.Permissions,
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

func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string, perms ...dto.APIKeyPerm) error {
	if ok := s.isValidAPIKeyFormat(apiKey); !ok {
		return errors.New(msg.ErrInvalidAPIKey)
	}

	apiKeyHash, err := s.kms.Hash(apiKey, s.apiKeySalt)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash api key")
		return errors.New(msg.ErrUnexpected)
	}

	apiKeyDetails, err := s.apiKeyRepo.GetApiKeyDetails(ctx, apiKeyHash)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errors.New(msg.ErrInvalidAPIKey)
		}
		return errors.New(msg.ErrUnexpected)
	}

	if apiKeyDetails.RevokedAt != nil && apiKeyDetails.RevokedAt.Before(time.Now()) {
		return errors.New(msg.ErrRevokedAPIKey)
	}

	if apiKeyDetails.ExpiresAt != nil && apiKeyDetails.ExpiresAt.Before(time.Now()) {
		return errors.New(msg.ErrExpiredAPIKey)
	}

	hasPerm := s.verifyAPIKeyPermissions(apiKeyDetails, perms...)

	if !hasPerm {
		return errors.New(msg.ErrInvalidAPIKey)
	}

	return nil
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
			if perm.ResouceID == apiPerm.ResourceID && perm.Perm == apiPerm.Permission {
				return true
			}
		}
	}
	return true
}
