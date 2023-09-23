package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New()

	server.Get("/hello", func(c *fiber.Ctx) error {
		c.SendString("Hello World, from Fiber!")
		return nil
	})

	server.Listen("0.0.0.0:8000")
}
