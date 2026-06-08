package middleware

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

type contextKey string

const (
	tenantIDKey   contextKey = "tenant_id"
	operatorIDKey contextKey = "operator_id"
	storeIDKey    contextKey = "store_id"
)

func ContextValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, tenantIDKey, parseInt64Header(c, "X-Tenant-ID", 1))
		ctx = context.WithValue(ctx, operatorIDKey, parseInt64Header(c, "X-Operator-ID", 0))
		ctx = context.WithValue(ctx, storeIDKey, parseInt64Header(c, "X-Store-ID", 0))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestIDFromContext(ctx context.Context) string {
	return requestid.FromContext(ctx)
}

func TenantIDFromContext(ctx context.Context) int64 {
	value, _ := ctx.Value(tenantIDKey).(int64)
	if value == 0 {
		return 1
	}
	return value
}

func OperatorIDFromContext(ctx context.Context) int64 {
	value, _ := ctx.Value(operatorIDKey).(int64)
	return value
}

func StoreIDFromContext(ctx context.Context) int64 {
	value, _ := ctx.Value(storeIDKey).(int64)
	return value
}

func parseInt64Header(c *gin.Context, key string, fallback int64) int64 {
	value := c.GetHeader(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}
