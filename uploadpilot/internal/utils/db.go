package utils

import (
	"go.mongodb.org/mongo-driver/bson"
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
