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

func ErrorFrom(c *gin.Context, err error) {
	ErrorFromWithDefault(c, err, "service error")
}

func ErrorFromWithDefault(c *gin.Context, err error, defaultMessage string) {
	appErr, ok := errno.As(err)
	if !ok {
		appErr = errno.New(errno.CodeInternalError, defaultMessage)
	}
	Error(c, HTTPStatus(appErr.Code), appErr.Code, appErr.Message, nil)
}

func HTTPStatus(code int) int {
	switch code {
	case errno.CodeInvalidParam:
		return http.StatusBadRequest
	case errno.CodeUnauthorized:
		return http.StatusUnauthorized
	case errno.CodeForbidden:
		return http.StatusForbidden
	case errno.CodeNotFound:
		return http.StatusNotFound
	case errno.CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
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
