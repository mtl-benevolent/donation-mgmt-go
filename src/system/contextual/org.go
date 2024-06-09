package contextual

import "context"

const OrgIdCtxKey ContextKey = "org_id"
const OrgSlugCtxKey ContextKey = "org_slug"

func WithOrgId(ctx context.Context, orgId int64) context.Context {
	return context.WithValue(ctx, OrgIdCtxKey, orgId)
}

func GetOrgId(ctx context.Context) int64 {
	orgId, ok := ctx.Value(OrgIdCtxKey).(int64)
	if !ok {
		return 0
	}

	return orgId
}

func WithOrgSlug(ctx context.Context, orgSlug string) context.Context {
	return context.WithValue(ctx, OrgSlugCtxKey, orgSlug)
}

func GetOrgSlug(ctx context.Context) string {
	orgId, ok := ctx.Value(OrgSlugCtxKey).(string)
	if !ok {
		return ""
	}

	return orgId
}
