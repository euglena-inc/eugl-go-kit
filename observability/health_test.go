package observability

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/euglena-inc/eugl-go-kit/requestid"
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

func TestReadinessDependencyStatusLogsDownAndRecovery(t *testing.T) {
	capture := &healthLogCapture{}
	health := NewHealth(ServiceInfo{ServiceName: "svc-test"}, nil, nil, slog.New(capture))
	ctx := requestid.WithContext(context.Background(), "rid-health-1")

	health.logDependencyStatus(ctx, "postgres", "down", errors.New("ping failed"))
	health.logDependencyStatus(ctx, "postgres", "up", nil)

	down := requireHealthLog(t, capture.records, "dependency readiness down")
	assertHealthLogString(t, down, "request_id", "rid-health-1")
	assertHealthLogString(t, down, "service_name", "svc-test")
	assertHealthLogString(t, down, "dependency", "postgres")
	assertHealthLogString(t, down, "status", "down")
	assertHealthLogString(t, down, "error", "ping failed")

	recovered := requireHealthLog(t, capture.records, "dependency readiness recovered")
	assertHealthLogString(t, recovered, "request_id", "rid-health-1")
	assertHealthLogString(t, recovered, "dependency", "postgres")
	assertHealthLogString(t, recovered, "status", "up")
}

func TestReadinessIncludesExtraDependencyChecks(t *testing.T) {
	capture := &healthLogCapture{}
	health := NewHealth(ServiceInfo{ServiceName: "svc-test"}, nil, nil, slog.New(capture))
	health.AddDependency("temporal", func(context.Context) error {
		return errors.New("temporal unavailable")
	})
	ctx := requestid.WithContext(context.Background(), "rid-health-2")

	data, ready := health.Readiness(ctx)
	if ready {
		t.Fatal("Readiness() ready = true, want false")
	}
	checks := data["checks"].(map[string]interface{})
	temporal := checks["temporal"].(map[string]interface{})
	if temporal["status"] != "down" {
		t.Fatalf("temporal status = %v, want down", temporal["status"])
	}

	down := requireHealthLog(t, capture.records, "dependency readiness down")
	assertHealthLogString(t, down, "request_id", "rid-health-2")
	assertHealthLogString(t, down, "dependency", "temporal")
}

func requireHealthLog(t *testing.T, records []map[string]slog.Value, message string) map[string]slog.Value {
	t.Helper()
	for _, record := range records {
		if record["message"].String() == message {
			return record
		}
	}
	t.Fatalf("log message %q missing in %+v", message, records)
	return nil
}

func assertHealthLogString(t *testing.T, record map[string]slog.Value, key string, want string) {
	t.Helper()
	got, ok := record[key]
	if !ok {
		t.Fatalf("log attr %q missing in %+v", key, record)
	}
	if got.String() != want {
		t.Fatalf("log attr %q = %q, want %q", key, got.String(), want)
	}
}
