package middlewares

import (
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

func PathParamsMiddleware(c *gin.Context) {
	if orgId, ok := c.Params.Get("orgId"); ok {
		c.Set(contextual.OrgIdCtxKey, orgId)
	}

	if orgSlug, ok := c.Params.Get("orgSlug"); ok {
		c.Set(contextual.OrgSlugCtxKey, orgSlug)
	}

	if env, ok := c.Params.Get("env"); ok {
		c.Set(contextual.EnvCtxKey, env)
	}
}
