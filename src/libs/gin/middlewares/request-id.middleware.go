package middlewares

import (
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIdHttpHeader = "X-Request-Id"

func RequestIdMiddleware(c *gin.Context) {
	requestId := c.Request.Header.Get(requestIdHttpHeader)
	if requestId == "" {
		requestId = uuid.NewString()
	}

	c.Set(string(contextual.RequestIdCtxKey), requestId)
	c.Header(requestIdHttpHeader, requestId)

	c.Next()
}
