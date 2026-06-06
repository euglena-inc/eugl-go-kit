package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if isHealthCheckPath(c.Request.URL.Path) {
			return
		}

		attrs := []slog.Attr{
			slog.String("request_id", requestid.FromContext(c.Request.Context())),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("http_status", c.Writer.Status()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
		}
		attrs = appendCommonBusinessIDs(c, attrs)

		log.LogAttrs(c.Request.Context(), slog.LevelInfo, "http_request", attrs...)
		if c.Writer.Status() >= 400 || len(c.Errors) > 0 {
			errorAttrs := append([]slog.Attr{}, attrs...)
			if errText := ginErrorText(c); errText != "" {
				errorAttrs = append(errorAttrs, slog.String("error", errText))
			}
			level := slog.LevelWarn
			if c.Writer.Status() >= 500 {
				level = slog.LevelError
			}
			log.LogAttrs(c.Request.Context(), level, "http_error", errorAttrs...)
		}
	}
}

func ginErrorText(c *gin.Context) string {
	if len(c.Errors) == 0 {
		return ""
	}
	if last := c.Errors.Last(); last != nil && last.Err != nil {
		return last.Err.Error()
	}
	return ""
}

func isHealthCheckPath(path string) bool {
	return path == "/healthz" || path == "/readyz"
}

func appendCommonBusinessIDs(c *gin.Context, attrs []slog.Attr) []slog.Attr {
	if brandID := firstNonEmpty(c.Param("brand_id"), c.Query("brand_id"), c.GetHeader("X-Brand-Id")); brandID != "" {
		attrs = append(attrs, slog.String("brand_id", brandID))
	}
	if storeID := firstNonEmpty(c.Param("store_id"), c.Query("store_id"), c.GetHeader("X-Store-Id")); storeID != "" {
		attrs = append(attrs, slog.String("store_id", storeID))
	}
	if orderNo := firstNonEmpty(c.Param("order_no"), c.Query("order_no")); orderNo != "" {
		attrs = append(attrs, slog.String("order_no", orderNo))
	}
	return attrs
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
