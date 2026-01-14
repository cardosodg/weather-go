package app

import (
	"WeatherTrack/internal/collector/config"
	"WeatherTrack/internal/collector/model"
	"WeatherTrack/internal/collector/service"
	"encoding/json"
	"log"
	"os"
	"sync"
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

func fetchHistory(locations []model.LocationInput) {
	log.Println("Querying history data...")
	time.Sleep(5 * time.Second)
	for _, loc := range locations {

		hist, _ := service.GetHistoryWeather(loc.Latitude, loc.Longitude, loc.Label)
		service.PostHistory(hist)
	}
}

func fetchWithRetry(loc model.LocationInput) (model.WeatherApiData, error) {

	var lastErr error

	for range 3 {
		data, err := service.GetSingleWeather(loc.Latitude, loc.Longitude, loc.Label)
		if err == nil {
			log.Printf("Data fetched for location %s", loc.Label)
			return data, nil
		}
		lastErr = err
		log.Printf("Failed to fetch data for %s. Retrying.", loc.Label)
		time.Sleep(1 * time.Second)
	}
	return model.WeatherApiData{}, lastErr
}

func sendWithRetry(data model.WeatherApiData) error {
	var lastErr error

	for range 5 {
		err := service.PostData(data)
		if err == nil {
			log.Printf("Data sent for location %s", data.Location)
			return nil
		}
		lastErr = err
		log.Printf("Failed to send data for location %s. Retrying.", data.Location)
		time.Sleep(300 * time.Millisecond)
	}

	return lastErr
}

func fetchSingleLocation(loc model.LocationInput) {
	data, err := fetchWithRetry(loc)
	if err != nil {
		log.Printf("WARN failed fetch for %s: %v", loc.Label, err)
		return
	}

	err = sendWithRetry(data)
	if err != nil {
		log.Printf("WARN failed send for %s: %v", loc.Label, err)
	}
}

func fetchCurrent(locations []model.LocationInput) {
	log.Println("Querying current data...")
	time.Sleep(5 * time.Second)

	for {
		var wg sync.WaitGroup
		wg.Add(len(locations))

		for _, loc := range locations {
			go func() {
				defer wg.Done()
				fetchSingleLocation(loc)
			}()

		}
		log.Println("Waiting for all locations to complete.")
		wg.Wait()
		log.Println("All data collected.")
		time.Sleep(config.OpenMeteoWaitTime * time.Minute)
	}
}

func Run() {
	locations := setupInit()

	fetchHistory(locations)
	fetchCurrent(locations)
}
