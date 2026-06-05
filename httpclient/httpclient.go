package httpclient

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

type Options struct {
	BaseURL string
	Timeout time.Duration
}

type Envelope[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	RequestID string `json:"request_id"`
}

func New(options Options) *resty.Client {
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	client := resty.New().
		SetBaseURL(options.BaseURL).
		SetTimeout(timeout).
		SetHeader("Content-Type", "application/json")

	client.OnBeforeRequest(func(_ *resty.Client, request *resty.Request) error {
		currentRequestID := requestid.FromContext(request.Context())
		if currentRequestID != "" {
			request.SetHeader(requestid.Header, currentRequestID)
		}
		return nil
	})

	return client
}

func UpstreamError(service string, messages ...string) error {
	for _, message := range messages {
		if message != "" {
			return fmt.Errorf("%s upstream error: %s", service, message)
		}
	}
	return fmt.Errorf("%s upstream error", service)
}
