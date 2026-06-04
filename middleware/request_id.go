package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRequestID := c.GetHeader(requestid.Header)
		if currentRequestID == "" {
			currentRequestID = requestid.New()
		}

		c.Header(requestid.Header, currentRequestID)
		c.Set("request_id", currentRequestID)
		c.Request = c.Request.WithContext(requestid.WithContext(c.Request.Context(), currentRequestID))
		c.Next()
	}
}
