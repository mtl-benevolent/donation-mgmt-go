package contextual

import (
	"context"

	"github.com/google/uuid"
)

const RequestIdCtxKey ContextKey = "request_id"

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, RequestIdCtxKey, requestId)
}

func GetRequestId(ctx context.Context) string {
	requestId, ok := ctx.Value(RequestIdCtxKey).(string)
	if !ok {
		return uuid.Nil.String()
	}

	return requestId
}
