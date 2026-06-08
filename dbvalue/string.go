package dbvalue

import "strings"

func NullableString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func StringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func NullableInt64(value int64) *int64 {
	if value <= 0 {
		return nil
	}
	return &value
}
