package main

import (
	"WeatherTrack/internal/collector/config"
	"WeatherTrack/internal/collector/model"
	"WeatherTrack/internal/collector/service"
	"encoding/json"
	"log"
	"os"
	"time"
)

func setupInit() []model.LocationInput {
	log.Println("Waiting 30 seconds before starting...")
	time.Sleep(30 * time.Second)

	data_raw, err := os.ReadFile(config.LocationsFile)
	if err != nil {
		log.Fatal(err)
	}

	var locations []model.LocationInput
	if err = json.Unmarshal(data_raw, &locations); err != nil {
		log.Fatal(err)
	}

	log.Println("The following locations will be monitored:")
	for _, loc := range locations {
		log.Println(loc.Label, loc.Latitude, loc.Longitude)
	}
	return locations
}

func main() {
	locations := setupInit()

	log.Println("Querying history data...")
	time.Sleep(5 * time.Second)
	for _, loc := range locations {

		hist, _ := service.GetHistoryWeather(loc.Latitude, loc.Longitude, loc.Label)
		service.PostHistory(hist)
	}

	log.Println("Querying current data...")
	time.Sleep(5 * time.Second)
	for {
		for _, loc := range locations {

			data, _ := service.GetSingleWeather(loc.Latitude, loc.Longitude, loc.Label)
			service.PostData(data)
		}
		time.Sleep(15 * time.Minute)
	}
}
