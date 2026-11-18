package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := newApp()
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(fmt.Errorf("gagal menjalankan server: %w", err))
	}
}
