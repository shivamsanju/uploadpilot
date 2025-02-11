package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Struct[T any] map[string]interface{}

func (s Struct[T]) Value() (driver.Value, error) {
	var temp T
	if err := mapstructure.Decode(s, &temp); err != nil {
		return nil, fmt.Errorf("invalid map structure: %w", err)
	}
	return json.Marshal(s)
}

func (s *Struct[T]) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan struct: invalid type")
	}

	return json.Unmarshal(b, s)
}
