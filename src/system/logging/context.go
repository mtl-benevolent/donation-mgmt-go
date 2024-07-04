package logging

import (
	"context"
	"log/slog"

	"donation-mgmt/src/system/contextual"
)

func WithContextData(ctx context.Context, logger *slog.Logger) *slog.Logger {
	return logger.With(ContextLogData(ctx)...)
}

func ContextLogData(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}

	attrs := make([]any, 0)
	attrs = append(attrs, slog.String(RequestIdKey, contextual.GetRequestId(ctx)))

	if orgId := contextual.GetOrgId(ctx); orgId != 0 {
		attrs = append(attrs, slog.Int64(OrgIdKey, orgId))
	}

	if env := contextual.GetEnv(ctx); env != "" {
		attrs = append(attrs, slog.String(EnvKey, env))
	}

	return attrs
}
