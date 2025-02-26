package dbutils

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	if errors.Is(err, gorm.ErrUnsupportedRelation) {
		return errors.New("unknown err")
	}

	return errors.New("database error: " + err.Error())
}

func GenerateUUID() string {
	return uuid.New().String()
}
