package middlewares

import (
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/logging"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type httpRequestLogEntry struct {
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	RequestSize   string `json:"requestSize,omitempty"`
	Status        int    `json:"status"`
	ResponseSize  string `json:"responseSize"`
	UserAgent     string `json:"userAgent"`
	Latency       string `json:"latency"`
	Protocol      string `json:"protocol"`
}

func LogRequestMiddleware(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	duration := time.Since(startTime)
	durationText := fmt.Sprintf("%.9fs", duration.Seconds())

	status := c.Writer.Status()

	logMessage := fmt.Sprintf("%s %d %s", c.Request.Method, status, c.Request.URL.String())

	httpReq := httpRequestLogEntry{
		RequestMethod: c.Request.Method,
		RequestURL:    c.Request.URL.String(),
		RequestSize:   c.Request.Header.Get("Content-Length"),
		Status:        status,
		ResponseSize:  fmt.Sprintf("%d", c.Writer.Size()),
		UserAgent:     c.Request.UserAgent(),
		Latency:       durationText,
		Protocol:      c.Request.Proto,
	}

	l := logger.ForComponent("WebServer").With(logging.ContextLogData(c)...)
	l.Info(logMessage, slog.String("handler", c.HandlerName()), slog.String("path", c.FullPath()), slog.Any("httpRequest", httpReq))
}
