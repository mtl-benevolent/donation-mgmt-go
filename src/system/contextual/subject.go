package contextual

import "context"

const SubjectCtxKey = "subject"

func WithSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, SubjectCtxKey, subject)
}

func GetSubject(ctx context.Context) string {
	subject, ok := ctx.Value(SubjectCtxKey).(string)
	if !ok {
		return ""
	}

	return subject
}
