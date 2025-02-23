package commonutils

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func IsInteger(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := strconv.Atoi(value)
	return err == nil
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

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("integer", IsInteger)
	validate.RegisterValidation("keyvaluepairs", IsKeyValuePairs)
	validate.RegisterValidation("sort", IsSortValid)
	return validate
}
