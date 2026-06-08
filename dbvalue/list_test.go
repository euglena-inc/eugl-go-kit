package dbvalue

import "testing"

func TestStringListValueAndScanJSON(t *testing.T) {
	value := StringList{"a.jpg", "b.jpg"}
	raw, err := value.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if raw != `["a.jpg","b.jpg"]` {
		t.Fatalf("Value() = %s", raw)
	}

	var scanned StringList
	if err := scanned.Scan(raw); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(scanned) != 2 || scanned[0] != "a.jpg" || scanned[1] != "b.jpg" {
		t.Fatalf("scanned = %#v", scanned)
	}
}

func TestStringListValueDefaultsToArray(t *testing.T) {
	raw, err := StringList(nil).Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if raw != "[]" {
		t.Fatalf("Value() = %s, want []", raw)
	}
}
