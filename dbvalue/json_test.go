package dbvalue

import (
	"database/sql/driver"
	"testing"
)

func TestJSONObjectTextValueDefaultsToObject(t *testing.T) {
	var value driver.Valuer = JSONObjectText("")
	got, err := value.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if got != "{}" {
		t.Fatalf("Value() = %#v, want object default", got)
	}
}

func TestJSONArrayTextValueDefaultsToArray(t *testing.T) {
	var value driver.Valuer = JSONArrayText("")
	got, err := value.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if got != "[]" {
		t.Fatalf("Value() = %#v, want array default", got)
	}
}

func TestJSONTextScanAcceptsDatabaseStringTypes(t *testing.T) {
	var object JSONObjectText
	if err := object.Scan([]byte(`{"a":1}`)); err != nil {
		t.Fatalf("Scan([]byte) error = %v", err)
	}
	if object != `{"a":1}` {
		t.Fatalf("object = %q", object)
	}

	var array JSONArrayText
	if err := array.Scan(`[{"weekday":1}]`); err != nil {
		t.Fatalf("Scan(string) error = %v", err)
	}
	if array != `[{"weekday":1}]` {
		t.Fatalf("array = %q", array)
	}
}

func TestJSONTextOrDefault(t *testing.T) {
	if got := JSONTextOrDefault("", "[]"); got != "[]" {
		t.Fatalf("JSONTextOrDefault(empty) = %q, want []", got)
	}
	if got := JSONTextOrDefault(" \t\n ", "{}"); got != "{}" {
		t.Fatalf("JSONTextOrDefault(blank) = %q, want {}", got)
	}
	if got := JSONTextOrDefault("{}", "[]"); got != "{}" {
		t.Fatalf("JSONTextOrDefault(value) = %q, want {}", got)
	}
}
