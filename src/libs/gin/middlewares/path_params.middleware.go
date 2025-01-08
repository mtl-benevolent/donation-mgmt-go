package middlewares

import (
	"donation-mgmt/src/libs/gin/ginext"
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

func PathParamsMiddleware(c *gin.Context) {
	if orgId, ok := c.Params.Get(ginext.OrgIDParamName); ok {
		c.Set(contextual.OrgIdCtxKey, orgId)
	}

	if orgSlug, ok := c.Params.Get(ginext.OrgSlugParamName); ok {
		c.Set(contextual.OrgSlugCtxKey, orgSlug)
	}

	if env, ok := c.Params.Get(ginext.EnvParamName); ok {
		c.Set(contextual.EnvCtxKey, env)
	}
}
