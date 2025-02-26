package repo

import (
	"context"
	"errors"

	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApiKeyRepo struct {
	db *driver.Driver
}

func NewApiKeyRepo(db *driver.Driver) *ApiKeyRepo {
	return &ApiKeyRepo{
		db: db,
	}
}

func (r *ApiKeyRepo) CreateApiKey(ctx context.Context, apiKey *models.APIKey) error {
	return r.db.Orm.WithContext(ctx).Create(apiKey).Error
}

func (r *ApiKeyRepo) GetApiKeyDetailsByID(ctx context.Context, id string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).First(&apiKey, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepo) GetApiKeyLimitedDetailsByID(ctx context.Context, id string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Preload("Permissions").
		Select("id", "name", "created_by", "created_at", "expires_at", "revoked", "revoked_at", "revoked_by", "can_read_acc", "can_manage_acc").
		First(&apiKey, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepo) GetApiKeyDetails(ctx context.Context, hash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Preload("Permissions").
		Preload("User").
		First(&apiKey, "api_key_hash = ?", hash).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepo) GetAllApiKeys(ctx context.Context, userID string) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Omit("user").
		Select("id", "name", "created_by", "created_at", "expires_at", "revoked", "revoked_at", "revoked_by").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "revoked"}, Desc: false},
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		}}).
		Find(&apiKeys, "user_id = ?", userID).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]models.APIKey, 0), nil
		}
		return nil, err
	}
	return apiKeys, nil
}

func (r *ApiKeyRepo) Update(ctx context.Context, apiKey *models.APIKey) error {
	return r.db.Orm.WithContext(ctx).Save(apiKey).Error
}
