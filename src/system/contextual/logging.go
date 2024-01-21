package contextual

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
	attrs = append(attrs, slog.String("requestId", GetRequestId(ctx)))

	return attrs
}
