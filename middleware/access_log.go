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
	}
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
