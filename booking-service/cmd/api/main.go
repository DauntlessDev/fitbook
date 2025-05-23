package main

import (
	"log"
	"os"

	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/config"
	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/server"
)

func main() {
	log.Println("Starting booking service...")

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	if err := server.Start(cfg); err != nil {
		log.Printf("Server exited with error: %v", err)
		os.Exit(1)
	}

	log.Println("Server exited cleanly")
}
