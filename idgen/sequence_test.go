package idgen

import (
	"testing"
	"time"
)

func TestNextIDWithClockUsesMillisecondPrefixAndSequenceSuffix(t *testing.T) {
	now := time.Date(2026, 6, 8, 10, 30, 0, 456*int(time.Millisecond), time.UTC)

	first := NextIDWithClock(func() time.Time { return now })
	second := NextIDWithClock(func() time.Time { return now })

	if first/10000 != now.UnixMilli() {
		t.Fatalf("prefix = %d, want %d", first/10000, now.UnixMilli())
	}
	if second != first+1 {
		t.Fatalf("second = %d, want first+1 %d", second, first+1)
	}
}
