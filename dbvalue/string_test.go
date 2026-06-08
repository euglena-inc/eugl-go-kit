package dbvalue

import "testing"

func TestNullableStringTrimsWhitespaceAndReturnsNilForEmpty(t *testing.T) {
	got := NullableString("  cashier  ")
	if got == nil || *got != "cashier" {
		t.Fatalf("NullableString() = %#v, want cashier pointer", got)
	}

	if empty := NullableString(" \t\n "); empty != nil {
		t.Fatalf("NullableString(empty) = %#v, want nil", *empty)
	}
}

func TestStringValueReturnsEmptyForNil(t *testing.T) {
	if got := StringValue(nil); got != "" {
		t.Fatalf("StringValue(nil) = %q, want empty string", got)
	}

	value := "manager"
	if got := StringValue(&value); got != value {
		t.Fatalf("StringValue() = %q, want %q", got, value)
	}
}
