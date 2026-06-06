package httpserver

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestRunShutsDownAndClosesOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	closed := false

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := Run(ctx, Options{
		Server: &http.Server{
			Addr:              "127.0.0.1:0",
			Handler:           http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
			ReadHeaderTimeout: time.Second,
		},
		ShutdownTimeout: time.Second,
		Close: func() error {
			closed = true
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !closed {
		t.Fatal("Close() was not called")
	}
}
