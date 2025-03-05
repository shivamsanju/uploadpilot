package rbac

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewGormAdapter(db *gorm.DB, tableName string) (*gormadapter.Adapter, error) {
	return gormadapter.NewAdapterByDBUseTableName(db, "", tableName)
}
