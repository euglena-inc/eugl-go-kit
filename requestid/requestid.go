package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

const Header = "X-Request-Id"

type contextKey string

const key contextKey = "request_id"

func New() string {
	var random [8]byte
	if _, err := rand.Read(random[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return time.Now().UTC().Format("20060102150405.000000000") + "-" + hex.EncodeToString(random[:])
}

func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, key, requestID)
}

func FromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(key).(string)
	return requestID
}
