package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed dist/*
var dist embed.FS

//go:embed dist/index.html
var indexHTML embed.FS

func RegisterHandlers(app *fiber.App) {
	if os.Getenv("ENV") == "dev" {
		log.Println("Running in dev mode")
		setupDevProxy(app)
		return
	}

	app.Get("/", func(c *fiber.Ctx) error {
		data, err := fs.ReadFile(indexHTML, "dist/index.html")
		if err != nil {
			return fiber.ErrInternalServerError
		}
		c.Set("Content-Type", "text/html")
		return c.Send(data)
	})

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(dist),
		PathPrefix: "dist",
		Browse:     false,
	}))

	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if shouldSkip(path) {
			return c.Next()
		}
		data, err := fs.ReadFile(indexHTML, "dist/index.html")
		if err != nil {
			return fiber.ErrInternalServerError
		}
		c.Set("Content-Type", "text/html")
		return c.Send(data)
	})
}

func setupDevProxy(app *fiber.App) {
	target := "http://localhost:5173"

	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if shouldSkip(path) {
			return c.Next()
		}
		return proxy.Do(c, target+c.OriginalURL())
	})
}

func shouldSkip(path string) bool {
	return len(path) >= 4 && path[:4] == "/api" ||
		len(path) >= 8 && path[:8] == "/swagger" ||
		len(path) >= 5 && path[:5] == "/auth"

}
