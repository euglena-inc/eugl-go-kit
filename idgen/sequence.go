package idgen

import (
	"sync/atomic"
	"time"
)

var defaultSequence uint64

func NextID() int64 {
	return NextIDWithClock(time.Now)
}

func NextIDWithClock(clock func() time.Time) int64 {
	if clock == nil {
		clock = time.Now
	}
	now := uint64(clock().UnixMilli())
	next := atomic.AddUint64(&defaultSequence, 1) % 10000
	return int64(now*10000 + next)
}
