package validation

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/uploadpilot/momentum/pkg/dsl"
)

func ValidateYamlStructure(workflow dsl.Workflow) error {
	var raw map[string]interface{}
	if err := mapstructure.Decode(workflow, &raw); err != nil {
		return err
	}
	log.Println(raw)
	if err := checkForExtraFields(raw, reflect.TypeOf(workflow)); err != nil {
		return err
	}
	fmt.Printf("Workflow: %+v\n", workflow)
	return validateStatement(workflow.Root)
}

func checkForExtraFields(data map[string]interface{}, expectedType reflect.Type) error {
	if expectedType.Kind() == reflect.Ptr {
		expectedType = expectedType.Elem()
	}

	// Ensure we are working with a struct
	if expectedType.Kind() != reflect.Struct {
		return nil
	}

	allowedFields := make(map[string]bool)
	fieldMap := make(map[string]reflect.StructField)

	for i := 0; i < expectedType.NumField(); i++ {
		field := expectedType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "-" {
			continue // Ignore fields marked as "-"
		}
		if yamlTag == "" {
			yamlTag = field.Name // Default to struct field name if no yaml tag
		}
		allowedFields[yamlTag] = true
		fieldMap[yamlTag] = field
	}

	log.Println(allowedFields)
	for key, value := range data {
		if !allowedFields[key] {
			return fmt.Errorf("unexpected field: %s", key)
		}
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if field, found := fieldMap[key]; found {
				if err := checkForExtraFields(nestedMap, field.Type); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func fieldMatchesYamlTag(fieldName, yamlTag string, t reflect.Type) bool {
	field, found := t.FieldByName(fieldName)
	if !found {
		return false
	}
	return field.Tag.Get("yaml") == yamlTag
}

func validateStatement(stmt dsl.Statement) error {
	fmt.Printf("S: %+v\n", stmt)

	count := 0
	if stmt.Activity != nil {
		count++
		if err := validateActivity(stmt.Activity); err != nil {
			return err
		}
	}
	if stmt.Sequence != nil {
		count++
		if err := validateSequence(stmt.Sequence); err != nil {
			return err
		}
	}
	if stmt.Parallel != nil {
		count++
		if err := validateParallel(stmt.Parallel); err != nil {
			return err
		}
	}
	fmt.Println("count: ", count)
	if count != 1 {
		return errors.New("a statement must contain exactly one of activity, sequence, or parallel")
	}
	return nil
}

func validateActivity(activity *dsl.ActivityInvocation) error {
	if activity.Name == "" {
		return errors.New("activity must have a name")
	}
	return nil
}

func validateSequence(seq *dsl.Sequence) error {
	if len(seq.Elements) == 0 {
		return errors.New("sequence must have at least one element")
	}
	for _, stmt := range seq.Elements {
		if err := validateStatement(*stmt); err != nil {
			return err
		}
	}
	return nil
}

func validateParallel(par *dsl.Parallel) error {
	if len(par.Branches) == 0 {
		return errors.New("parallel must have at least one branch")
	}
	for _, stmt := range par.Branches {
		if err := validateStatement(*stmt); err != nil {
			return err
		}
	}
	return nil
}
