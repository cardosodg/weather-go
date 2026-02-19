package service

import (
	"WeatherTrack/internal/forecast/config"
	receiverModel "WeatherTrack/internal/receiver/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthCheck struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func GetHealth() (HealthCheck, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	var healthCheck HealthCheck

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

func PostForecast(data []receiverModel.WeatherDTO) error {

	if len(data) == 0 {
		return nil
	}

	client := &http.Client{Timeout: 10 * time.Second}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		config.ForecastURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
