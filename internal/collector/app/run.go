package app

import (
	"WeatherTrack/internal/collector/config"
	"WeatherTrack/internal/collector/service"
	receiverModel "WeatherTrack/internal/receiver/model"
	"sync"

	"encoding/json"
	"log"
	"os"
	"time"
)

type LocationInput struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Label     string `json:"label"`
}

func setupInit() []LocationInput {
	log.Println("Starting collector...")
	time.Sleep(5 * time.Second)

	dataRaw, err := os.ReadFile(config.LocationsFile)
	if err != nil {
		log.Fatal(err)
	}

	var locations []LocationInput
	if err = json.Unmarshal(dataRaw, &locations); err != nil {
		log.Fatal(err)
	}

	log.Println("The following locations will be monitored:")
	for _, loc := range locations {
		log.Println(loc.Label, loc.Latitude, loc.Longitude)
	}

	return locations
}

func fetchHistory(locations []LocationInput) {
	log.Println("Starting history data fetch and transmission...")
	time.Sleep(5 * time.Second)

	for _, loc := range locations {
		hist, err := service.GetHistoryWeather(loc.Latitude, loc.Longitude, loc.Label)
		if err != nil {
			log.Printf("Failed to fetch history for %s: %v", loc.Label, err)
			continue
		}

		if err := service.PostHistory(hist); err != nil {
			log.Printf("Failed to send history for %s: %v", loc.Label, err)
		}
	}
	log.Println("History data fetch and transmission completed.")
}

func fetchWithRetry(loc LocationInput) (receiverModel.WeatherDTO, error) {

	var lastErr error
	var dto receiverModel.WeatherDTO

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

	return dto, lastErr
}

func sendWithRetry(data receiverModel.WeatherDTO) error {

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

func fetchSingleLocation(loc LocationInput) {
	data, err := fetchWithRetry(loc)
	if err != nil {
		log.Printf("WARN failed fetch for %s: %v", loc.Label, err)
		return
	}

	if err := sendWithRetry(data); err != nil {
		log.Printf("WARN failed send for %s: %v", loc.Label, err)
	}
}

func fetchCurrent(locations []LocationInput) {
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

		time.Sleep(config.CollectorWaitTime * time.Minute)
	}
}

func checkReceiver() bool {
	for i := 1; i <= config.ReceiverMaxRetries; i++ {
		healthCheck, err := service.GetHealth()
		if err != nil {
			log.Printf("health check attempt %d/%d error: %v",
				i, config.ReceiverMaxRetries, err)
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
		config.ReceiverMaxRetries)

	return false
}

func Run() {
	locations := setupInit()

	if checkReceiver() {
		fetchHistory(locations)
		fetchCurrent(locations)
	}
}
