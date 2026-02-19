package app

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/service"
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type LocationInput struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Label     string `json:"label"`
}

func setupInit() []LocationInput {
	log.Println("Starting forecast...")
	time.Sleep(5 * time.Second)

	dataRaw, err := os.ReadFile(config.LocationsFile)
	if err != nil {
		log.Fatal(err)
	}

	var locations []LocationInput
	if err = json.Unmarshal(dataRaw, &locations); err != nil {
		log.Fatal(err)
	}

	log.Println("Generating forecast for:")
	for _, loc := range locations {
		log.Println(loc.Label, loc.Latitude, loc.Longitude)
	}

	return locations
}

func fetchWithRetry(loc LocationInput) ([]receiverModel.WeatherDTO, error) {
	var lastErr error

	for range 3 {
		data, err := service.GetForecastWeather(
			loc.Latitude,
			loc.Longitude,
			loc.Label,
		)

		if err == nil {
			log.Printf("Forecast fetched for %s (%d items)", loc.Label, len(data))
			return data, nil
		}

		lastErr = err
		log.Printf("Failed to fetch forecast for %s. Retrying...", loc.Label)
		time.Sleep(1 * time.Second)
	}

	return nil, lastErr
}

func sendWithRetry(batch []receiverModel.WeatherDTO) error {
	var lastErr error

	for range 5 {
		err := service.PostForecast(batch)
		if err == nil {
			log.Printf("Forecast batch sent for location %s (%d items)",
				batch[0].Location,
				len(batch),
			)
			return nil
		}

		lastErr = err
		log.Printf("Failed to send data for location %s. Retrying.", batch[0].Location)
		time.Sleep(300 * time.Millisecond)
	}

	return lastErr
}

func fetchSingleForecastLocation(loc LocationInput) {
	batch, err := fetchWithRetry(loc)
	if err != nil {
		log.Printf("WARN fetch failed for %s: %v", loc.Label, err)
		return
	}

	if len(batch) == 0 {
		log.Printf("No forecast data for %s", loc.Label)
		return
	}

	err = sendWithRetry(batch)
	if err != nil {
		log.Printf("WARN send failed for %s: %v", loc.Label, err)
	}
}

func fetchForecast(locations []LocationInput) {
	log.Println("Querying forecast data...")
	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup
	wg.Add(len(locations))

	for _, loc := range locations {
		loc := loc // evita problema de closure
		go func() {
			defer wg.Done()
			fetchSingleForecastLocation(loc)
		}()
	}

	log.Println("Waiting for all locations to complete.")
	wg.Wait()
	log.Println("All forecast data collected.")
}

func checkReceiver() bool {
	for i := 1; i <= config.ReceiverMaxRetries; i++ {
		healthCheck, err := service.GetHealth()
		if err != nil {
			log.Printf("health check attempt %d/%d error: %v",
				i,
				config.ReceiverMaxRetries,
				err,
			)
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

	log.Printf("receiver did not become healthy after %d attempts",
		config.ReceiverMaxRetries,
	)
	return false
}

func Run() {
	locations := setupInit()

	if checkReceiver() {
		fetchForecast(locations)
	}
}
