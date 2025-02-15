package commonutils

import (
	"fmt"
	"reflect"
)

func ConvertDTOToModel[T any, M any](dto *T, model *M) error {
	if dto == nil || model == nil {
		return fmt.Errorf("dto and model must not be nil")
	}

	dtoVal := reflect.ValueOf(dto).Elem()     // Dereference DTO pointer
	modelVal := reflect.ValueOf(model).Elem() // Dereference Model pointer

	if !modelVal.CanAddr() {
		return fmt.Errorf("model must be a pointer to a struct")
	}

	dtoType := dtoVal.Type()
	for i := 0; i < dtoVal.NumField(); i++ {
		field := dtoType.Field(i)
		dtoFieldValue := dtoVal.Field(i)

		modelField := modelVal.FieldByName(field.Name)

		// Check if the field exists in Model and can be set
		if modelField.IsValid() && modelField.CanSet() && modelField.Type() == dtoFieldValue.Type() {
			modelField.Set(dtoFieldValue)
		}
	}
	return nil
}
