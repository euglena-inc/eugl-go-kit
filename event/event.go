package event

import (
	"context"
	"time"
)

type Event struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	EventVersion  int                    `json:"event_version"`
	OccurredAt    time.Time              `json:"occurred_at"`
	SourceService string                 `json:"source_service"`
	RequestID     string                 `json:"request_id"`
	TraceID       string                 `json:"trace_id"`
	BrandID       string                 `json:"brand_id,omitempty"`
	StoreID       string                 `json:"store_id,omitempty"`
	Payload       map[string]interface{} `json:"payload"`
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(_ context.Context, _ Event) error {
	return nil
}
