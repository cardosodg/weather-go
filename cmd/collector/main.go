package main

import "WeatherTrack/internal/collector/service"

func main() {
	lat := "-20.317865"
	lon := "-40.310211"
	service.GetSingleWeather("-20.3222", "-40.3381", "vitoria")
	service.GetSingleWeather(lat, lon, "vitoria")
}
