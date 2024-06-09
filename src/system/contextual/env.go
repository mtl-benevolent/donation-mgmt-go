package contextual

import "context"

const EnvCtxKey = "env"

func WithEnv(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, EnvCtxKey, env)
}

func GetEnv(ctx context.Context) string {
	env, ok := ctx.Value(EnvCtxKey).(string)
	if !ok {
		return ""
	}

	return env
}
