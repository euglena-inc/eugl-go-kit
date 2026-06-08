package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/errno"
	"github.com/euglena-inc/eugl-go-kit/requestid"
	"github.com/euglena-inc/eugl-go-kit/response"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered",
			slog.String("request_id", requestid.FromContext(c.Request.Context())),
			slog.Any("panic", recovered),
		)
		response.Error(c, http.StatusInternalServerError, errno.CodeInternalError, "service error", nil)
	})
}
