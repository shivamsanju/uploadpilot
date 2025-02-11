package validations

import (
	"fmt"

	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateJSONSchema(schema string, json map[string]interface{}) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	val, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		infra.Log.Errorf("error loading schema: %s", err.Error())
		return fmt.Errorf("error loading schema")
	}

	jsonDoc := gojsonschema.NewGoLoader(json)
	result, err := val.Validate(jsonDoc)
	if err != nil {
		infra.Log.Errorf("error validating JSON: %s", err.Error())
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
