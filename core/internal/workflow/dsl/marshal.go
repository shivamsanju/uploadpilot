package dsl

import "encoding/json"

func (w Workflow) MarshalJSON() ([]byte, error) {
	type Alias Workflow
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&w),
	})
}

func (w Workflow) MarshalYAML() (interface{}, error) {
	type Alias Workflow
	return &struct {
		*Alias
	}{
		Alias: (*Alias)(&w),
	}, nil
}
