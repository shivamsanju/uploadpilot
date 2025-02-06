package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"gorm.io/gorm"
)

func DBError(err error) error {
	if err == nil {
		return nil
	}
	infra.Log.Errorf("DATABASE ERROR: %s", err.Error())

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}
	return errors.New("there was an issue processing your request. please try again later")
}

func GenerateUUID() string {
	return uuid.New().String()
}
