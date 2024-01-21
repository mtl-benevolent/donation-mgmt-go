package middlewares

import (
	"donation-mgmt/src/system/context_data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIdHttpHeader = "X-Request-Id"

func RequestIdMiddleware(c *gin.Context) {
	requestId := c.Request.Header.Get(requestIdHttpHeader)
	if requestId == "" {
		requestId = uuid.NewString()
	}

	c.Set(context_data.RequestIdCtxKey, requestId)
	c.Header(requestIdHttpHeader, requestId)

	c.Next()
}
