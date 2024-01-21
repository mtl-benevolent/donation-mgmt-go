package main

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/gin"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/organizations"
	"fmt"
	"log/slog"
	"os"

	"github.com/gretro/go-lifecycle"
)

func main() {
	appConfig := config.Bootstrap()
	logger := logger.BootstrapLogger(appConfig)

	gs := lifecycle.NewGracefulShutdown(context.Background())
	readyCheck := lifecycle.NewReadyCheck()

	defer func() {
		err := recover()
		if err != nil {
			logger.Error("Application panicked", slog.Any("error", err))
			fmt.Printf("Application panicked: %v\n", err)

			os.Exit(1)
		}
	}()

	db.Bootstrap(gs, readyCheck, appConfig)

	router := gin.Bootstrap(gs, readyCheck, appConfig)

	organizations.Bootstrap(router)

	readyCheck.StartPolling()
	logger.Info("Application is ready")

	gs.WaitForShutdown()
	logger.Info("Application is shutting down")
}
