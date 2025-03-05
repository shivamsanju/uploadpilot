package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	validate.RegisterValidation("integer", IsInteger)
	validate.RegisterValidation("future", IsFutureTime)
	validate.RegisterValidation("keyvaluepairs", IsKeyValuePairs)
	validate.RegisterValidation("sort", IsSortValid)
	validate.RegisterValidation("alphanumspace", IsAlphaNumSpace)

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) ValidateStruct(s any) error {
	if err := v.validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("%s (%s%s)", e.Field(), e.Tag(), e.Param()))
			}
			return errors.New(strings.Join(errMsgs, "; "))
		}
		return err
	}
	return nil
}

func (v *Validator) ValidateJSONSchema(schema string, json map[string]interface{}) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	val, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		zap.L().Error(fmt.Sprintf("VALIDATION_ERROR : error loading schema: %s", err.Error()))
		return fmt.Errorf("error loading schema")
	}

	jsonDoc := gojsonschema.NewGoLoader(json)
	result, err := val.Validate(jsonDoc)
	if err != nil {
		zap.L().Error(fmt.Sprintf("VALIDATION_ERROR : error validating JSON: %s", err.Error()))
		return err
	}

	if !result.Valid() {
		e := ""
		for _, desc := range result.Errors() {
			e += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("JSON is not valid: %s", e)
	}

	return nil
}
