package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// generateMockGeoJSON returns a mock FeatureCollection for 3D map rendering
func generateMockGeoJSON() fiber.Map {
	return fiber.Map{
		"type": "FeatureCollection",
		"features": []fiber.Map{
			{
				"type": "Feature",
				"properties": fiber.Map{
					"id":       "bld-001",
					"name":     "Student Dormitory Block A",
					"type":     "residential",
					"height":   15,
					"capacity": 120,
				},
				"geometry": fiber.Map{
					"type": "Polygon",
					"coordinates": [][][]float64{
						{
							{95.3697, 5.5683},
							{95.3699, 5.5683},
							{95.3699, 5.5681},
							{95.3697, 5.5681},
							{95.3697, 5.5683},
						},
					},
				},
			},
			{
				"type": "Feature",
				"properties": fiber.Map{
					"id":       "bld-002",
					"name":     "Main Library",
					"type":     "academic",
					"height":   24,
					"capacity": 500,
				},
				"geometry": fiber.Map{
					"type": "Polygon",
					"coordinates": [][][]float64{
						{
							{95.3710, 5.5690},
							{95.3715, 5.5690},
							{95.3715, 5.5685},
							{95.3710, 5.5685},
							{95.3710, 5.5690},
						},
					},
				},
			},
		},
	}
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// API Route Group
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	// Get Spatial Data (GeoJSON)
	api.Get("/spatial/buildings", func(c *fiber.Ctx) error {
		geoData := generateMockGeoJSON()

		// Optional: Implement simple category filtering if query param exists
		categoryFilter := c.Query("type")
		if categoryFilter != "" {
			var filteredFeatures []fiber.Map
			features := geoData["features"].([]fiber.Map)

			for _, feature := range features {
				props := feature["properties"].(fiber.Map)
				if props["type"] == categoryFilter {
					filteredFeatures = append(filteredFeatures, feature)
				}
			}
			geoData["features"] = filteredFeatures
		}

		return c.JSON(geoData)
	})

	log.Fatal(app.Listen(":8080"))
}
