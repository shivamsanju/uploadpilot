package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/kms"
)

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	infra.Log.Infof("scanning JSONB... value: %+v", value)
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (j *JSONB) String() string {
	data, _ := json.MarshalIndent(j, "", "  ") // Pretty print
	return string(data)
}

type EncryptedJSONB map[string]interface{}

func (b EncryptedJSONB) Value() (driver.Value, error) {
	v, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	if len(string(v)) == 0 {
		return nil, nil
	}

	ev, err := kms.Encrypt(config.EncryptionKey, string(v))
	if err != nil {
		return nil, err
	}
	return ev, nil
}

func (b *EncryptedJSONB) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(v) == 0 {
		return nil
	}
	dv, err := kms.Decrypt(config.EncryptionKey, v)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(dv), &b)
}

func (j *EncryptedJSONB) String() string {
	data, _ := json.MarshalIndent(j, "", "  ") // Pretty print
	return string(data)
}
