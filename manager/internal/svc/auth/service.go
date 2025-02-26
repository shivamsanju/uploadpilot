package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/phuslu/log"
	"github.com/quasoft/memstore"
	commonutils "github.com/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/go-core/common/vault"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/manager/internal/utils"
)

var jwtSecretBytes []byte

func InitJWTSessionStore(jwtSecret string) error {
	store := memstore.NewMemStore(
		[]byte("authkey123"),
	)
	gothic.Store = store

	jwtSecret = strings.TrimSpace(jwtSecret)
	if len(jwtSecret) < 32 {
		return errors.New("JWT secret must be at least 32 characters long")
	}

	jwtSecretBytes = []byte(jwtSecret)

	goth.UseProviders(
		google.New(config.GoogleClientID, config.GoogleClientSecret, config.GoogleCallbackURL, "email", "profile"),
		github.New(config.GithubClientID, config.GithubClientSecret, config.GithubCallbackURL, "email", "user"),
	)
	return nil
}

type Service struct {
	apiKeyRepo *repo.ApiKeyRepo
	kms        vault.KMS
	apiKeySalt string
}

func NewService(apiKeyRepo *repo.ApiKeyRepo, kms vault.KMS) *Service {
	return &Service{
		apiKeyRepo: apiKeyRepo,
		kms:        kms,
		apiKeySalt: config.ApiKeyEncryptionSalt,
	}
}

func (s *Service) GetAllApiKeysForUser(ctx context.Context) ([]models.APIKey, error) {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.apiKeyRepo.GetAllApiKeys(ctx, user.UserID)
}

func (s *Service) CreateApiKey(ctx context.Context, data *dto.CreateApiKeyData) (string, error) {
	user, err := utils.GetUserDetailsFromContext(ctx)
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
		Name:       data.Name,
		UserID:     user.UserID,
		ApiKeyHash: hashedKey,
		ExpiresAt:  data.ExpiresAt,
		Revoked:    false,
		CreatedByColumn: models.CreatedByColumn{
			CreatedBy: user.Email,
		},
	}

	if err := s.apiKeyRepo.CreateApiKey(ctx, apiKey); err != nil {
		return "", err
	}

	return newKey, nil
}

func (s *Service) GetApiKeyDetails(ctx context.Context, key string) (*models.APIKey, error) {
	return s.apiKeyRepo.GetApiKeyDetails(ctx, key)
}

func (s *Service) RevokeApiKey(ctx context.Context, id string) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	apiKey, err := s.apiKeyRepo.GetApiKeyDetailsByID(ctx, id)
	if err != nil {
		return err
	}
	apiKey.Revoked = true
	apiKey.RevokedAt = time.Now()
	apiKey.RevokedBy = user.UserID
	return s.apiKeyRepo.Update(ctx, apiKey)
}

func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string) (*dto.APIKeyClaims, error) {
	if ok := s.isValidAPIKeyFormat(apiKey); !ok {
		return nil, errors.New(msg.ErrInvalidAPIKey)
	}

	apiKeyHash, err := s.kms.Hash(apiKey, s.apiKeySalt)
	if err != nil {
		return nil, err
	}

	apiKeyDetails, err := s.apiKeyRepo.GetApiKeyDetails(ctx, apiKeyHash)
	if err != nil {
		return nil, err
	}

	if apiKeyDetails.Revoked || apiKeyDetails.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return &dto.APIKeyClaims{
		UserClaims: &dto.UserClaims{
			UserID: apiKeyDetails.UserID,
			Email:  apiKeyDetails.User.Email,
			Name:   *apiKeyDetails.User.Name,
		},
		CanReadAcc:   apiKeyDetails.CanReadAcc,
		CanManageAcc: apiKeyDetails.CanManageAcc,
		Permissions:  apiKeyDetails.Permissions,
	}, nil
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

func (s *Service) GenerateJWTToken(w http.ResponseWriter, user *models.User, expiry time.Duration) (string, error) {
	if user.Name == nil {
		name := "User"
		user.Name = &name
	}
	claims := &dto.JWTClaims{
		UserClaims: &dto.UserClaims{
			UserID: user.ID,
			Email:  user.Email,
			Name:   *user.Name,
		},
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(expiry).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecretBytes)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ValidateJWTToken(jwtToken string) (*dto.JWTClaims, error) {
	parts := strings.Split(jwtToken, " ")

	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) == 0 {
		return nil, errors.New(msg.InvalidToken)
	}

	token, err := jwt.ParseWithClaims(
		parts[1],
		&dto.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecretBytes, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.JWTClaims)
	if !ok {
		return claims, errors.New(msg.InvalidToken)
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return claims, errors.New(msg.TokenExpired)
	}
	return claims, nil
}
