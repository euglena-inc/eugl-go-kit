package dbvalue

import (
	"database/sql/driver"
	"fmt"
)

type JSONObjectText string

func (j JSONObjectText) Value() (driver.Value, error) {
	if j == "" {
		return "{}", nil
	}
	return string(j), nil
}

func (j *JSONObjectText) Scan(value any) error {
	return scanJSONText(value, func(text string) {
		*j = JSONObjectText(text)
	})
}

type JSONArrayText string

func (j JSONArrayText) Value() (driver.Value, error) {
	if j == "" {
		return "[]", nil
	}
	return string(j), nil
}

func (j *JSONArrayText) Scan(value any) error {
	return scanJSONText(value, func(text string) {
		*j = JSONArrayText(text)
	})
}

func scanJSONText(value any, assign func(string)) error {
	switch v := value.(type) {
	case nil:
		assign("")
	case []byte:
		assign(string(v))
	case string:
		assign(v)
	default:
		return fmt.Errorf("scan json text: %T", value)
	}
	return nil
}
