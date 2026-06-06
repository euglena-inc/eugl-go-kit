package observability

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

type healthLogCapture struct {
	records []map[string]slog.Value
}

func (h *healthLogCapture) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *healthLogCapture) Handle(_ context.Context, record slog.Record) error {
	attrs := map[string]slog.Value{"message": slog.StringValue(record.Message)}
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.records = append(h.records, attrs)
	return nil
}

func (h *healthLogCapture) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *healthLogCapture) WithGroup(string) slog.Handler {
	return h
}

func TestReadinessDoesNotLogDependencyStatus(t *testing.T) {
	capture := &healthLogCapture{}
	health := NewHealth(ServiceInfo{ServiceName: "svc-test"}, nil, nil, slog.New(capture))
	health.AddDependency("temporal", func(context.Context) error {
		return errors.New("temporal unavailable")
	})

	_, _ = health.Readiness(context.Background())

	if len(capture.records) != 0 {
		t.Fatalf("Readiness() logs = %+v, want none", capture.records)
	}
}

func TestReadinessIncludesExtraDependencyChecks(t *testing.T) {
	capture := &healthLogCapture{}
	health := NewHealth(ServiceInfo{ServiceName: "svc-test"}, nil, nil, slog.New(capture))
	health.AddDependency("temporal", func(context.Context) error {
		return errors.New("temporal unavailable")
	})
	data, ready := health.Readiness(context.Background())
	if ready {
		t.Fatal("Readiness() ready = true, want false")
	}
	checks := data["checks"].(map[string]interface{})
	temporal := checks["temporal"].(map[string]interface{})
	if temporal["status"] != "down" {
		t.Fatalf("temporal status = %v, want down", temporal["status"])
	}
}
