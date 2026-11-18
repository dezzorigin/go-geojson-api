package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// version adalah versi API yang dilaporkan di route root ("/").
const version = "1.1.0"

// newApp membangun aplikasi Fiber dengan semua route.
// Dipisah dari main() supaya bisa diuji dengan app.Test tanpa
// menjalankan server sungguhan.
func newApp() *fiber.App {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Indeks API: daftar endpoint yang tersedia.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "go-geojson-api",
			"version": version,
			"endpoints": []string{
				"GET /api/v1/health",
				"GET /api/v1/spatial/buildings",
				"GET /api/v1/spatial/buildings?type=residential",
				"GET /api/v1/spatial/buildings?type=academic",
			},
		})
	})

	api := app.Group("/api/v1")

	// Health check.
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	// Data spasial (GeoJSON).
	api.Get("/spatial/buildings", listBuildings)

	return app
}
