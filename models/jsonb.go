package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONB map[string]interface{}

func JSONBFromObject(i interface{}) (JSONB, error) {
	j := JSONB{}
	b, err := json.Marshal(i)
	if err != nil {
		return j, err
	}

	if err := json.Unmarshal(b, &j); err != nil {
		return j, err
	}
	return j, nil
}

func (j *JSONB) Bytes() []byte {
	b, err := json.Marshal(j)
	if err != nil {
		return nil
	}
	return b
}

func (j *JSONB) WithBytes(b []byte) {
	if j == nil {
		return
	}

	err := json.Unmarshal(b, &j)
	if err != nil {
		return
	}
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, j)

	}
	return errors.New(fmt.Sprint("Failed to unmarshal JSON from DB", src))
}
