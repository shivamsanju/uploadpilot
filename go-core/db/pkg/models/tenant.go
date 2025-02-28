package models

type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusInActive TenantStatus = "inactive"
)

type Tenant struct {
	ID           string                 `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name         string                 `gorm:"column:name;not null;type:varchar(50)" json:"name"`
	OwnerID      string                 `gorm:"column:owner_id;not null;type:uuid" json:"ownerId"` // Reference to the primary user/owner
	ContactEmail string                 `gorm:"column:contact_email;not null;type:varchar(100)" json:"contactEmail"`
	Phone        string                 `gorm:"column:phone;type:varchar(20)" json:"phone,omitempty"`
	Address      string                 `gorm:"column:address;type:text" json:"address,omitempty"`
	Industry     string                 `gorm:"column:industry;type:varchar(50)" json:"industry,omitempty"`
	Status       TenantStatus           `gorm:"column:status;not null;type:varchar(20);default:'active'" json:"status"` // e.g., active, suspended
	Metadata     map[string]interface{} `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAtColumn
	UpdatedAtColumn
}

func (Tenant) TableName() string {
	return "tenant"
}
