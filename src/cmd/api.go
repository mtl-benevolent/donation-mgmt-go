package main

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/fiber"
	"fmt"
	"os"

	"github.com/gretro/go-lifecycle"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("Application panicked\n")

			if appErr, ok := err.(error); ok {
				fmt.Printf("Panic cause: " + appErr.Error())
			}

			os.Exit(1)
		}
	}()

	gs := lifecycle.NewGracefulShutdown(context.Background())

	appConfig := config.Bootstrap()

	fiber.Bootstrap(gs, appConfig)

	gs.WaitForShutdown()

	fmt.Println("Application is shutting down")
}
