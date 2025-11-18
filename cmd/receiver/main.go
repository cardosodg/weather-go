package main

import (
	"WeatherTrack/internal/receiver/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("receiver")
	host := "0.0.0.0"
	port := 8123
	address := fmt.Sprintf("%s:%d", host, port)

	r := gin.Default()

	router.SetupRoutes(r)

	err := r.Run(address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
