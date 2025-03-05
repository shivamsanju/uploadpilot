package validator

import (
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func IsInteger(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := strconv.Atoi(value)
	return err == nil
}

func IsFutureTime(fl validator.FieldLevel) bool {
	dateTime := fl.Field().Interface().(time.Time)
	now := time.Now().UTC()
	return dateTime.After(now)
}

// IsKeyValuePairs validates search query (field:value,field:value,...)
func IsKeyValuePairs(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Allow empty values
	}

	// Updated regex: key:value1,value2;key2:value3,value4
	pattern := `^([a-zA-Z0-9_]+:[^;]+)(;[a-zA-Z0-9_]+:[^;]+)*$`
	match, _ := regexp.MatchString(pattern, value)
	return match
}

func IsSortValid(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Allow empty values
	}
	pattern := `^[a-zA-Z0-9_]+:(asc|desc)$`
	match, _ := regexp.MatchString(pattern, value)
	return match
}

func IsAlphaNumSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Allow empty values
	}
	pattern := `^[a-zA-Z0-9 ]+$`
	match, _ := regexp.MatchString(pattern, value)
	return match
}
