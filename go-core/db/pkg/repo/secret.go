package repo

import (
	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/models"
)

type SecretRepo struct {
	db *driver.Driver
}

func NewSecretRepo(db *driver.Driver) *SecretRepo {
	return &SecretRepo{
		db: db,
	}
}

func (s *SecretRepo) GetAllSecretsWithValues(workspaceID string) ([]models.Secret, error) {
	var secrets []models.Secret
	err := s.db.Orm.Find(&secrets, "workspace_id = ?", workspaceID).Error
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *SecretRepo) GetAllSecretsWithoutValues(workspaceID string) ([]models.Secret, error) {
	var secrets []models.Secret
	err := s.db.Orm.Omit("value").Omit("salt").Find(&secrets, "workspace_id = ?", workspaceID).Error
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *SecretRepo) GetSecretWithValue(workspaceID, key string) (*models.Secret, error) {
	var secret models.Secret
	err := s.db.Orm.First(&secret, "workspace_id = ? AND key = ?", workspaceID, key).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func (s *SecretRepo) GetSecretWithoutValue(workspaceID, key string) (*models.Secret, error) {
	var secret models.Secret
	err := s.db.Orm.Omit("value").Omit("salt").First(&secret, "workspace_id = ? AND key = ?", workspaceID, key).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func (s *SecretRepo) CreateSecret(secret *models.Secret) error {
	return s.db.Orm.Create(secret).Error
}

func (s *SecretRepo) BulkCreateSecrets(secrets []models.Secret) error {
	return s.db.Orm.Create(secrets).Error
}

func (s *SecretRepo) UpdateSecret(secret *models.Secret) error {
	return s.db.Orm.Save(secret).Error
}

func (s *SecretRepo) BulkUpdateSecrets(secrets []models.Secret) error {
	return s.db.Orm.Save(secrets).Error
}

func (s *SecretRepo) DeleteSecret(workspaceID, key string) error {
	return s.db.Orm.Delete(&models.Secret{}, "workspace_id = ? AND key = ?", workspaceID, key).Error
}

func (s *SecretRepo) BulkDeleteSecrets(workspaceID string, keys []string) error {
	return s.db.Orm.Delete(&models.Secret{}, "workspace_id = ? AND key IN (?)", workspaceID, keys).Error
}

func (s *SecretRepo) DeleteAllSecrets(workspaceID string) error {
	return s.db.Orm.Delete(&models.Secret{}, "workspace_id = ?", workspaceID).Error
}
