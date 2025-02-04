package utils

import (
	"errors"

	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

func FilterNonEmptyBsonFields(data bson.M) bson.M {
	filtered := bson.M{}
	for key, value := range data {
		if value != nil && value != "" {
			filtered[key] = value
		}
	}
	return filtered
}

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
