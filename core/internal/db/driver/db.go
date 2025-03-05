package driver

import "gorm.io/gorm"

type Driver struct {
	Orm *gorm.DB
}
