package db

import "gorm.io/gorm"

type DB struct {
	Orm *gorm.DB
}
