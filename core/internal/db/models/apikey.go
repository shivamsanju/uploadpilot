package models

import (
	"time"

	"github.com/uploadpilot/core/internal/db/dtypes"
)

type APIKey struct {
	ID          string             `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string             `gorm:"column:name;not null;type:varchar(50)" json:"name"`
	TenantID    string             `gorm:"column:tenant_id;not null;type:uuid" json:"tenantId"`
	ApiKeyHash  string             `gorm:"column:api_key_hash;not null;type:varchar(128)" json:"-"`
	LastUsedAt  *time.Time         `gorm:"column:last_used_at" json:"lastUsedAt,omitempty"`
	ExpiresAt   *time.Time         `gorm:"column:expires_at" json:"expiresAt,omitempty"`
	RevokedAt   *time.Time         `gorm:"column:revoked_at" json:"revokedAt,omitempty"`
	RevokedBy   *string            `gorm:"column:revoked_by;type:uuid" json:"revokedBy,omitempty"`
	UserID      string             `gorm:"column:user_id;not null;type:text" json:"userId"`
	IPWhitelist dtypes.StringArray `gorm:"column:ip_whitelist;type:text[]" json:"ipWhitelist,omitempty"`
	CreatedAtColumn
	UpdatedAtColumn
}

// TableName overrides the table name for GORM
func (APIKey) TableName() string {
	return "api_keys"
}
