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

func TestNullableInt64ReturnsNilForNonPositive(t *testing.T) {
	for _, value := range []int64{0, -1} {
		if got := NullableInt64(value); got != nil {
			t.Fatalf("NullableInt64(%d) = %#v, want nil", value, got)
		}
	}

	got := NullableInt64(42)
	if got == nil || *got != 42 {
		t.Fatalf("NullableInt64(42) = %#v, want 42 pointer", got)
	}
}
