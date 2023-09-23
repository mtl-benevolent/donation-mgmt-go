package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New()

	server.Get("/hello", func(c *fiber.Ctx) error {
		err := c.SendString("Hello World, from Fiber!")
		return err
	})

	err := server.Listen("0.0.0.0:8000")
	if err != nil {
		panic("Could not start server: " + err.Error())
	}
}
