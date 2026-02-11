package service

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func GetHealth() (model.HealthCheck, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	var healthCheck model.HealthCheck

	resp, err := client.Get(config.HealthURL)
	if err != nil {
		return healthCheck, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return healthCheck, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&healthCheck); err != nil {
		return healthCheck, err
	}

	return healthCheck, nil

}

func PostData(data model.WeatherApiForecast) error {
	client := &http.Client{Timeout: 5 * time.Second}
	body, _ := json.Marshal(data)
	log.Printf("Posting data: %s", string(body))

	resp, err := client.Post(
		config.ForecastURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
