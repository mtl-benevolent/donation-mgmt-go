package contextual

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"fmt"
	"strings"
)

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

func GetValidEnv(ctx context.Context) (dal.Environment, error) {
	env := GetEnv(ctx)

	switch strings.ToUpper(env) {
	case string(dal.EnvironmentSANDBOX):
		return dal.EnvironmentSANDBOX, nil
	case string(dal.EnvironmentLIVE):
		return dal.EnvironmentLIVE, nil
	default:
		return "", &apperrors.ValidationError{
			EntityName: "Environment",
			InnerError: fmt.Errorf("invalid environment value in path"),
		}
	}
}
