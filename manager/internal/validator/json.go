package validator

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

func ValidateJSONSchema(schema string, json map[string]interface{}) error {
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
