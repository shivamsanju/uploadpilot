package dtypes

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan value into StringArray: %v", value)
	}

	v = v[1 : len(v)-1]

	if len(v) == 0 {
		*s = StringArray{}
		return nil
	}

	*s = strings.Split(v, ",")
	return nil
}

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	// Convert the StringArray to the PostgreSQL array format
	return "{" + strings.Join(s, ",") + "}", nil
}

func (s StringArray) ArrayValue() ([]string, error) {
	return s, nil
}
