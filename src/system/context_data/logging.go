package context_data

import (
	"context"
	"log/slog"
)

func LoggerWithContextData(ctx context.Context, logger *slog.Logger) *slog.Logger {
	return logger.With(ContextLogData(ctx)...)
}

func ContextLogData(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}

	attrs := make([]any, 0)
	attrs = append(attrs, slog.String("request_id", GetRequestId(ctx)))

	return attrs
}
