package main

import (
	"WeatherTrack/internal/receiver/config"
	"WeatherTrack/internal/receiver/database"
	"WeatherTrack/internal/receiver/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting receiver...")

	db, err := database.Initialize()

	if err != nil {
		log.Print(err)
	}

	address := fmt.Sprintf("%s:%d", config.ReceiverIP, config.ReceiverPort)

	r := gin.Default()

	router.SetupRoutes(r, db)

	err = r.Run(address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
