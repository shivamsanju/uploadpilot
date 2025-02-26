package models

import "time"

type APIKey struct {
	ID           string       `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name         string       `gorm:"column:name;not null;type:varchar(25)" json:"name"`
	ApiKeyHash   string       `gorm:"column:api_key_hash;not null;type:varchar(128)" json:"apiKeyHash"`
	LastUsedAt   time.Time    `gorm:"column:last_used_at;" json:"lastUsedAt"`
	ExpiresAt    time.Time    `gorm:"column:expires_at;not null" json:"expiresAt"`
	Revoked      bool         `gorm:"column:revoked;not null;default:false" json:"revoked"`
	RevokedAt    time.Time    `gorm:"column:revoked_at" json:"revokedAt"`
	RevokedBy    string       `gorm:"column:revoked_by" json:"revokedBy"`
	UserID       string       `gorm:"column:user_id;not null;type:uuid" json:"userId"`
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	CanReadAcc   bool         `gorm:"column:can_read_acc;not null;default:false" json:"canReadAcc"`
	CanManageAcc bool         `gorm:"column:can_manage_acc;not null;default:false" json:"canManageAcc"`
	Permissions  []APIKeyPerm `gorm:"foreignKey:APIKeyID;constraint:OnDelete:CASCADE" json:"permissions"`
	CreatedByColumn
	CreatedAtColumn
}

func (APIKey) TableName() string {
	return "api_keys"
}

type APIKeyPerm struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	APIKeyID    string    `gorm:"column:api_key_id;not null;type:uuid" json:"apiKeyId"`
	WorkspaceID string    `gorm:"column:workspace_id;not null;type:uuid" json:"workspaceId"`
	CanRead     bool      `gorm:"column:can_read;not null;default:false" json:"canRead"`
	CanManage   bool      `gorm:"column:can_manage;not null;default:false" json:"canManage"`
	CanUpload   bool      `gorm:"column:can_upload;not null;default:false" json:"canUpload"`
	APIKey      APIKey    `gorm:"foreignKey:APIKeyID"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID"`
}

func (APIKeyPerm) TableName() string {
	return "api_key_permissions"
}
