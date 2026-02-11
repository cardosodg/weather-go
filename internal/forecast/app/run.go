package app

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/model"
	"WeatherTrack/internal/forecast/service"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

func setupInit() []model.LocationInput {
	log.Println("Starting forecast...")
	time.Sleep(5 * time.Second)

	data_raw, err := os.ReadFile(config.LocationsFile)
	if err != nil {
		log.Fatal(err)
	}

	var locations []model.LocationInput
	if err = json.Unmarshal(data_raw, &locations); err != nil {
		log.Fatal(err)
	}

	log.Println("Generating forecast for:")
	for _, loc := range locations {
		log.Println(loc.Label, loc.Latitude, loc.Longitude)
	}
	return locations
}

func fetchWithRetry(loc model.LocationInput) (model.WeatherApiForecast, error) {

	var lastErr error

	for range 3 {
		data, err := service.GetForecastWeather(loc.Latitude, loc.Longitude, loc.Label)
		if err == nil {
			log.Printf("Data fetched for location %s", loc.Label)
			return data, nil
		}
		lastErr = err
		log.Printf("Failed to fetch data for %s. Retrying.", loc.Label)
		time.Sleep(1 * time.Second)
	}
	return model.WeatherApiForecast{}, lastErr
}

func sendWithRetry(data model.WeatherApiForecast) error {
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

func fetchSingleForecastLocation(loc model.LocationInput) {
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

func fetchForecast(locations []model.LocationInput) {
	log.Println("Querying forecast data...")
	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup
	wg.Add(len(locations))

	for _, loc := range locations {
		go func() {
			defer wg.Done()
			fetchSingleForecastLocation(loc)
		}()

	}
	log.Println("Waiting for all locations to complete.")
	wg.Wait()
	log.Println("All data collected.")
}

func checkReceiver() bool {

	for i := 1; i <= config.ReceiverMaxRetries; i++ {
		healthCheck, err := service.GetHealth()
		if err != nil {
			log.Printf("health check attempt %d/%d error: %v", i, config.ReceiverMaxRetries, err)
			time.Sleep(config.ReceiverRetryInterval)
			continue
		}

		if healthCheck.Status == "ok" {
			log.Printf("receiver healthy on attempt %d", i)
			return true
		}

		log.Printf("receiver not ready (status=%s) attempt %d/%d",
			healthCheck.Status,
			i,
			config.ReceiverMaxRetries,
		)

		time.Sleep(config.ReceiverRetryInterval)
	}

	log.Printf("receiver did not become healthy after %d attempts", config.ReceiverMaxRetries)
	return false
}

func Run() {
	locations := setupInit()
	if checkReceiver() {
		fetchForecast(locations)
	}
}
