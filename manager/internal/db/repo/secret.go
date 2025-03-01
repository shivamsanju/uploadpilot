package repo

import (
	"context"

	"github.com/uploadpilot/manager/internal/db/driver"
	"github.com/uploadpilot/manager/internal/db/models"
	dbutils "github.com/uploadpilot/manager/internal/db/utils"
)

type SecretRepo struct {
	db *driver.Driver
}

func NewSecretRepo(db *driver.Driver) *SecretRepo {
	return &SecretRepo{
		db: db,
	}
}

func (r *SecretRepo) GetAllSecretsWithValues(ctx context.Context, workspaceID string) ([]models.Secret, error) {
	var secrets []models.Secret
	err := r.db.Orm.WithContext(ctx).Find(&secrets, "workspace_id = ?", workspaceID).Error
	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return secrets, nil
}

func (r *SecretRepo) GetAllSecretsWithoutValues(ctx context.Context, workspaceID string) ([]models.Secret, error) {
	var secrets []models.Secret
	err := r.db.Orm.WithContext(ctx).Omit("value").Omit("salt").Find(&secrets, "workspace_id = ?", workspaceID).Error
	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return secrets, nil
}

func (r *SecretRepo) GetSecretWithValue(ctx context.Context, workspaceID, key string) (*models.Secret, error) {
	var secret models.Secret
	err := r.db.Orm.WithContext(ctx).First(&secret, "workspace_id = ? AND key = ?", workspaceID, key).Error
	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &secret, nil
}

func (r *SecretRepo) GetSecretWithoutValue(ctx context.Context, workspaceID, key string) (*models.Secret, error) {
	var secret models.Secret
	err := r.db.Orm.WithContext(ctx).Omit("value").Omit("salt").First(&secret, "workspace_id = ? AND key = ?", workspaceID, key).Error
	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &secret, nil
}

func (r *SecretRepo) CreateSecret(ctx context.Context, secret *models.Secret) error {
	if err := r.db.Orm.WithContext(ctx).Create(secret).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SecretRepo) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	if err := r.db.Orm.WithContext(ctx).Save(secret).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SecretRepo) BulkUpdateSecrets(ctx context.Context, secrets []models.Secret) error {
	if err := r.db.Orm.WithContext(ctx).Save(secrets).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SecretRepo) DeleteSecret(ctx context.Context, workspaceID, key string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Secret{}, "workspace_id = ? AND key = ?", workspaceID, key).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SecretRepo) BulkDeleteSecrets(ctx context.Context, workspaceID string, keys []string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Secret{}, "workspace_id = ? AND key IN (?)", workspaceID, keys).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SecretRepo) DeleteAllSecrets(ctx context.Context, workspaceID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Secret{}, "workspace_id = ?", workspaceID).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
