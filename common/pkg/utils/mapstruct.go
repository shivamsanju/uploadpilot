package commonutils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func MapStructAndValidate(m map[string]interface{}, result interface{}) error {
	err := mapstructure.Decode(m, result)
	if err != nil {
		return fmt.Errorf("map decoding failed: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(result)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
