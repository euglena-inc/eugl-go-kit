package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/euglena-inc/eugl-go-kit/requestid"
)

type Options struct {
	BaseURL string
	Timeout time.Duration
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
