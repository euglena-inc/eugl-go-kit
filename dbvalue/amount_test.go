package dbvalue

import "testing"

func TestAmountCentsValueFormatsDecimal(t *testing.T) {
	raw, err := NewAmountCents(-12345).Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if raw != "-123.45" {
		t.Fatalf("Value() = %s, want -123.45", raw)
	}
}

func TestAmountCentsScanDecimalString(t *testing.T) {
	var amount AmountCents
	if err := amount.Scan("19.9"); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if amount.Cents() != 1990 {
		t.Fatalf("Cents() = %d, want 1990", amount.Cents())
	}
}

func TestAmountCentsScanFloatRoundsToCents(t *testing.T) {
	var amount AmountCents
	if err := amount.Scan(12.345); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if amount.Cents() != 1235 {
		t.Fatalf("Cents() = %d, want 1235", amount.Cents())
	}
}
