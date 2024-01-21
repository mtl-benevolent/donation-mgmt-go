package gin

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/libs/logger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gretro/go-lifecycle"
)

var router *gin.Engine

func Bootstrap(gs *lifecycle.GracefulShutdown, rc *lifecycle.ReadyCheck, appConfig *config.AppConfiguration) *gin.Engine {
	l := logger.ForComponent("Gin")

	router = gin.Default()
	if appConfig.AppEnvironment != config.Development {
		gin.SetMode(gin.ReleaseMode)
	}

	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "Healthy!")
	})

	router.GET("/ready", func(ctx *gin.Context) {
		if rc.Ready() {
			ctx.String(200, "Ready!")
		} else {
			ctx.String(500, "Not ready")
		}
	})

	router.GET("/ready/explain", func(ctx *gin.Context) {
		ctx.JSON(200, rc.Explain())
	})

	router.Use(middlewares.RequestIdMiddleware)
	router.Use(gin.Logger())

	router.Use(middlewares.ErrorHandler)
	router.Use(gin.CustomRecovery(middlewares.PanicHandler))

	// TODO: Implement error handling middleware

	router.Any("/panic", func(ctx *gin.Context) {
		panic("Test Panic")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", appConfig.HTTPPort),
		Handler: router,
	}

	go func() {
		l.Info("Starting Web Server", slog.String("addr", server.Addr))
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			l.Error("Web server shut down unexpectedly", slog.Any("error", err))

			_ = gs.Shutdown()
			os.Exit(1)
		}
	}()

	gs.RegisterComponentWithFn("WebServer", func() error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		return err
	})

	return router
}
