package fiber

import (
	"donation-mgmt/src/config"
	"fmt"
	"time"

	fiberlib "github.com/gofiber/fiber/v2"
	"github.com/gretro/go-lifecycle"
)

const shutdownTimeout = 2 * time.Second

var server *fiberlib.App

func Bootstrap(gs *lifecycle.GracefulShutdown, appConfig *config.AppConfiguration) *fiberlib.App {
	server = fiberlib.New()

	server.Get("/hello", func(c *fiberlib.Ctx) error {
		err := c.SendString("Hello World, from Fiber!")
		return err
	})

	go func() {
		err := server.Listen(fmt.Sprintf("0.0.0.0:%d", appConfig.HTTPPort))
		if err != nil {
			panic("Could not start server: " + err.Error())
		}
	}()

	gs.RegisterComponentWithFn("Fiber Server", func() error {
		err := server.ShutdownWithTimeout(shutdownTimeout)
		return err
	})

	return server
}

func Server() *fiberlib.App {
	if server != nil {
		panic("Fiber App has not been bootstrapped")
	}

	return server
}
