package main

import (
	"donation-mgmt/src/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	appConfig := config.Bootstrap()

	server := fiber.New()

	server.Get("/hello", func(c *fiber.Ctx) error {
		err := c.SendString("Hello World, from Fiber!")
		return err
	})

	err := server.Listen(fmt.Sprintf("0.0.0.0:%d", appConfig.HttpPort))
	if err != nil {
		panic("Could not start server: " + err.Error())
	}
}
