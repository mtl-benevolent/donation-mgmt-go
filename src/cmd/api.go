package main

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/fiber"
	"fmt"
	"os"

	"github.com/gretro/go-lifecycle"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("Application panicked: %v\n", err)

			os.Exit(1)
		}
	}()

	gs := lifecycle.NewGracefulShutdown(context.Background())

	appConfig := config.Bootstrap()

	db.Bootstrap(gs, appConfig)
	fiber.Bootstrap(gs, appConfig)

	gs.WaitForShutdown()

	fmt.Println("Application is shutting down")
}
