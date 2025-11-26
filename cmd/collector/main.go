package main

import (
	"WeatherTrack/internal/collector/service"
	"log"
	"time"
)

func main() {
	lat := "-20.317865"
	lon := "-40.310211"
	//data, _ := service.GetSingleWeather(lat, lon, "vitoria")
	// log.Println(data)
	// service.PostData(data)

	log.Println("Querying history data...")
	time.Sleep(5 * time.Second)
	hist, _ := service.GetHistoryWeather(lat, lon, "vitoria")
	service.PostHistory(hist)
}
