package idgen

import (
	"bytes"
	"testing"
	"time"
)

func TestTimeGeneratorUsesMillisecondPrefixAndRandomSuffix(t *testing.T) {
	now := time.Date(2026, 6, 8, 9, 30, 15, 123*int(time.Millisecond), time.UTC)
	generator := NewTimeGenerator(
		WithClock(func() time.Time { return now }),
		WithReader(bytes.NewReader([]byte{0x12, 0x34})),
	)

	got := generator.NextID()
	want := now.UnixMilli()*10000 + 0x1234%10000

	if got != want {
		t.Fatalf("NextID() = %d, want %d", got, want)
	}
}

func TestTimeGeneratorFallsBackToZeroSuffixWhenReaderFails(t *testing.T) {
	now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	generator := NewTimeGenerator(
		WithClock(func() time.Time { return now }),
		WithReader(bytes.NewReader(nil)),
	)

	got := generator.NextID()
	want := now.UnixMilli() * 10000

	if got != want {
		t.Fatalf("NextID() = %d, want %d", got, want)
	}
}
