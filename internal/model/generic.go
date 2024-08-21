package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonbArray[T any] []T

func (a *JsonbArray[T]) Scan(raw interface{}) error {
	if raw == nil {
		*a = nil
		return nil
	}
	switch v := raw.(type) {
	case []byte:
		bytes := raw.([]byte)
		if len(bytes) == 0 {
			return nil
		}
		return json.Unmarshal(bytes, &a)
	default:
		return fmt.Errorf("cannot sql.Scan() from: %#v", v)
	}
}

func (a JsonbArray[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	res, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("cannot driver.Value() from: %#v", a)
	}
	return driver.Value(res), nil
}
