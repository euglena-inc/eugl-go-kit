package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

type captureHandler struct {
	records []map[string]slog.Value
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := map[string]slog.Value{"message": slog.StringValue(record.Message)}
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.records = append(h.records, attrs)
	return nil
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := &captureHandler{}
	for _, attr := range attrs {
		next.records = append(next.records, map[string]slog.Value{attr.Key: attr.Value})
	}
	return next
}

func (h *captureHandler) WithGroup(string) slog.Handler {
	return h
}

func TestAccessLogIncludesRequestAndBusinessIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	capture := &captureHandler{}
	router := gin.New()
	router.Use(RequestID(), AccessLog(slog.New(capture)))
	router.GET("/api/v1/stores/:store_id/orders/:order_no", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stores/2001/orders/ORD123?brand_id=1001", nil)
	req.Header.Set(requestid.Header, "rid-1")
	router.ServeHTTP(httptest.NewRecorder(), req)

	record := requireLog(t, capture.records, "http_request")
	assertLogString(t, record, "request_id", "rid-1")
	assertLogString(t, record, "brand_id", "1001")
	assertLogString(t, record, "store_id", "2001")
	assertLogString(t, record, "order_no", "ORD123")
}

func TestAccessLogSkipsHealthChecks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	capture := &captureHandler{}
	router := gin.New()
	router.Use(RequestID(), AccessLog(slog.New(capture)))
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/readyz", func(c *gin.Context) { c.Status(http.StatusServiceUnavailable) })

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if len(capture.records) != 0 {
		t.Fatalf("health check logs = %+v, want none", capture.records)
	}
}

func TestAccessLogWritesErrorLogWithBusinessIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	capture := &captureHandler{}
	router := gin.New()
	router.Use(RequestID(), AccessLog(slog.New(capture)))
	router.POST("/api/v1/stores/:store_id/orders/:order_no", func(c *gin.Context) {
		_ = c.Error(context.Canceled)
		c.Status(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/stores/2001/orders/ORD123?brand_id=1001", nil)
	req.Header.Set(requestid.Header, "rid-err-1")
	router.ServeHTTP(httptest.NewRecorder(), req)

	record := requireLog(t, capture.records, "http_error")
	assertLogString(t, record, "request_id", "rid-err-1")
	assertLogString(t, record, "brand_id", "1001")
	assertLogString(t, record, "store_id", "2001")
	assertLogString(t, record, "order_no", "ORD123")
	assertLogString(t, record, "error", "context canceled")
}

func requireLog(t *testing.T, records []map[string]slog.Value, message string) map[string]slog.Value {
	t.Helper()
	for _, record := range records {
		if record["message"].String() == message {
			return record
		}
	}
	t.Fatalf("log message %q missing in %+v", message, records)
	return nil
}

func assertLogString(t *testing.T, attrs map[string]slog.Value, key string, want string) {
	t.Helper()
	got, ok := attrs[key]
	if !ok {
		t.Fatalf("log attr %q missing in %+v", key, attrs)
	}
	if got.String() != want {
		t.Fatalf("log attr %q = %q, want %q", key, got.String(), want)
	}
}
