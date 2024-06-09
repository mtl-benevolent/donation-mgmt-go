package middlewares

import (
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

const DevHeader = "x-user"
const DefaultSubject = "root"

func DevHeadersAuthMiddleware(c *gin.Context) {
	subject := c.GetHeader(DevHeader)
	if subject == "" {
		subject = DefaultSubject
	}

	c.Set(contextual.SubjectCtxKey, subject)
	c.Next()
}
