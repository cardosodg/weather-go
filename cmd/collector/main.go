package main

import (
	"WeatherTrack/internal/collector/service"
	"log"
)

func main() {
	lat := "-20.317865"
	lon := "-40.310211"
	data, _ := service.GetSingleWeather("-20.3222", "-40.3381", "vitoria")
	log.Println(data)
	service.PostData(data)
	service.GetSingleWeather(lat, lon, "vitoria")
}
