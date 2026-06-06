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
	attrs map[string]slog.Value
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	if h.attrs == nil {
		h.attrs = map[string]slog.Value{}
	}
	record.Attrs(func(attr slog.Attr) bool {
		h.attrs[attr.Key] = attr.Value
		return true
	})
	return nil
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := &captureHandler{attrs: map[string]slog.Value{}}
	for key, value := range h.attrs {
		next.attrs[key] = value
	}
	for _, attr := range attrs {
		next.attrs[attr.Key] = attr.Value
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

	assertLogString(t, capture.attrs, "request_id", "rid-1")
	assertLogString(t, capture.attrs, "brand_id", "1001")
	assertLogString(t, capture.attrs, "store_id", "2001")
	assertLogString(t, capture.attrs, "order_no", "ORD123")
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
