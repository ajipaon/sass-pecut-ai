package main

import (
	"sass-with-ai/web"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	web.RegisterHandlers(app)

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
