package main

import (
	"github.com/gofiber/fiber/v2"
)

// building adalah satu gedung dalam data contoh.
// Koordinat memakai [longitude, latitude] sesuai spesifikasi GeoJSON
// (RFC 7946). Nilai di bawah adalah data contoh di kawasan Darussalam,
// Banda Aceh — bukan hasil survei.
//
// NOTE: data ini mock statis di dalam kode. API ini tidak punya database
// dan tidak membaca file eksternal; endpoint selalu mengembalikan dua
// gedung yang sama di bawah ini.
func buildings() []fiber.Map {
	return []fiber.Map{
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
	}
}

// listBuildings mengembalikan FeatureCollection GeoJSON berisi gedung.
// Query param opsional `type` memfilter berdasarkan properties.type
// (mis. ?type=residential). Filter yang tidak cocok dengan gedung mana
// pun mengembalikan `"features": []` — array kosong, bukan null — supaya
// respons tetap GeoJSON yang valid.
func listBuildings(c *fiber.Ctx) error {
	features := buildings()

	if filter := c.Query("type"); filter != "" {
		filtered := []fiber.Map{}
		for _, feature := range features {
			props := feature["properties"].(fiber.Map)
			if props["type"] == filter {
				filtered = append(filtered, feature)
			}
		}
		features = filtered
	}

	return c.JSON(fiber.Map{
		"type":     "FeatureCollection",
		"features": features,
	})
}
