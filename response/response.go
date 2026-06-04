package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/errno"
	"github.com/euglena-inc/eugl-go-kit/requestid"
)

type Envelope struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, errno.CodeSuccess, "success", data)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, errno.CodeSuccess, "success", data)
}

func Error(c *gin.Context, statusCode int, code int, message string, data interface{}) {
	JSON(c, statusCode, code, message, data)
}

func JSON(c *gin.Context, statusCode int, code int, message string, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

	c.JSON(statusCode, Envelope{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: requestid.FromContext(c.Request.Context()),
	})
}
