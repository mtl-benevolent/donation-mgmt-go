package main

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/donations"
	"donation-mgmt/src/libs/db"
	firebaseadmin "donation-mgmt/src/libs/firebase-admin"
	"donation-mgmt/src/libs/gin"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/organizations"
	"donation-mgmt/src/permissions"
	"log/slog"
	"os"

	"github.com/gretro/go-lifecycle"
)

func main() {
	appConfig := config.Bootstrap()
	logger := logger.BootstrapLogger(appConfig)

	appConfig.WarnUnsafeOptions(logger)

	gs := lifecycle.NewGracefulShutdown(context.Background())
	readyCheck := lifecycle.NewReadyCheck()

	defer func() {
		err := recover()
		if err != nil {
			logger.Error("Application panicked", slog.Any("error", err))

			os.Exit(1)
		}
	}()

	if appConfig.EnableFirebase() {
		logger.Info("Firebase services are required. Bootstrapping client...")
		firebaseadmin.Bootstrap(appConfig)
	}

	db.Bootstrap(gs, readyCheck, appConfig)

	router := gin.Bootstrap(gs, readyCheck, appConfig)

	permissions.Bootstrap()
	organizations.Bootstrap(router)
	donations.Bootstrap(router)

	readyCheck.StartPolling()
	logger.Info("Application is ready")

	if err := gs.WaitForShutdown(); err != nil {
		panic(err)
	}

	logger.Info("Application is shutting down")
}
