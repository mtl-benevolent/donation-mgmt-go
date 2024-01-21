package gin

import (
	"donation-mgmt/src/config"

	"github.com/gin-gonic/gin"
	"github.com/gretro/go-lifecycle"
)

var server *gin.Engine

func Bootstrap(gs *lifecycle.GracefulShutdown, appConfig *config.AppConfiguration) {
	server = gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

}
