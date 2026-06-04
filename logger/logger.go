package logger

import (
	"log/slog"
	"os"
)

type Options struct {
	ServiceName string
	Env         string
	Level       slog.Level
}

func New(options Options) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: options.Level,
	})

	attrs := []slog.Attr{}
	if options.ServiceName != "" {
		attrs = append(attrs, slog.String("service_name", options.ServiceName))
	}
	if options.Env != "" {
		attrs = append(attrs, slog.String("env", options.Env))
	}

	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}

	return slog.New(handler).With(args...)
}
