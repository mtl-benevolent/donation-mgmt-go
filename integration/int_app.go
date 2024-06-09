package integration

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gretro/go-lifecycle"

	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/organizations"
)

type IntegrationApp struct {
	gs         *lifecycle.GracefulShutdown
	readyCheck *lifecycle.ReadyCheck
}

func NewIntegrationApp() *IntegrationApp {
	return &IntegrationApp{
		gs:         lifecycle.NewGracefulShutdown(context.Background()),
		readyCheck: lifecycle.NewReadyCheck(),
	}
}

func (app *IntegrationApp) Start(ctx context.Context) error {
	os.Setenv("APP_NAME", "int-tests")

	cfg := config.Bootstrap()
	logger.BootstrapLogger(cfg)

	db.Bootstrap(app.gs, app.readyCheck, cfg)

	// Init modules
	organizations.Bootstrap(nil)

	app.readyCheck.StartPolling()

	err := app.WaitForReady(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for app to be ready: %w", err)
	}

	return nil
}

func (app *IntegrationApp) WaitForReady(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if app.readyCheck.Ready() {
				fmt.Println("App is ready")
				return nil
			}
		}

	}
}

func (app *IntegrationApp) Stop() error {
	err := app.gs.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	return nil
}
