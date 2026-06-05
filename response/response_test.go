package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/euglena-inc/eugl-go-kit/errno"
)

func TestHTTPStatusMapsErrnoCodes(t *testing.T) {
	tests := []struct {
		code int
		want int
	}{
		{code: errno.CodeInvalidParam, want: http.StatusBadRequest},
		{code: errno.CodeUnauthorized, want: http.StatusUnauthorized},
		{code: errno.CodeForbidden, want: http.StatusForbidden},
		{code: errno.CodeNotFound, want: http.StatusNotFound},
		{code: errno.CodeConflict, want: http.StatusConflict},
		{code: errno.CodeInternalError, want: http.StatusInternalServerError},
		{code: 9999, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		if got := HTTPStatus(tt.code); got != tt.want {
			t.Fatalf("HTTPStatus(%d) = %d, want %d", tt.code, got, tt.want)
		}
	}
}

func TestErrorFromWritesErrnoEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/stores", nil)

	ErrorFrom(c, errno.New(errno.CodeNotFound, "store not found"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	if got := rec.Body.String(); got != `{"code":1004,"message":"store not found","data":{},"request_id":""}` {
		t.Fatalf("body = %s", got)
	}
}

func TestErrorFromMasksUnknownErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/stores", nil)

	ErrorFrom(c, errors.New("database password leaked"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if got := rec.Body.String(); got != `{"code":1006,"message":"service error","data":{},"request_id":""}` {
		t.Fatalf("body = %s", got)
	}
}
