package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

// JSON is a custom type that wraps a map to hold JSON objects.
type JSON map[string]interface{}

func (j *JSON) SetBSON(raw bson.Raw) error {
	var data interface{}
	if err := bson.Unmarshal(raw, &data); err != nil {
		return err
	}
	*j = data.(map[string]interface{})
	return nil
}

func (j JSON) GetBSON() (interface{}, error) {
	return bson.Marshal(j)
}

func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(j))
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*j = obj
	return nil
}

func (j JSON) String() string {
	data, _ := json.MarshalIndent(j, "", "  ") // Pretty print
	return string(data)
}
