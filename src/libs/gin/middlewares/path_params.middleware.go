package middlewares

import (
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

func PathParamsMiddleware(c *gin.Context) {
	if orgId, ok := c.Params.Get("orgId"); ok {
		c.Set(string(contextual.OrgIdCtxKey), orgId)
	}

	if orgSlug, ok := c.Params.Get("orgSlug"); ok {
		c.Set(string(contextual.OrgSlugCtxKey), orgSlug)
	}

	if env, ok := c.Params.Get("env"); ok {
		c.Set(string(contextual.EnvCtxKey), env)
	}
}
