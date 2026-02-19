package service

import (
	"WeatherTrack/internal/collector/config"
	receiverModel "WeatherTrack/internal/receiver/model"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type healthResponse struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func GetHealth() (healthResponse, error) {

	client := &http.Client{Timeout: 5 * time.Second}
	var health healthResponse

	resp, err := client.Get(config.HealthURL)
	if err != nil {
		return health, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return health, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return health, err
	}

	return health, nil
}

func PostData(data receiverModel.WeatherDTO) error {

	client := &http.Client{Timeout: 5 * time.Second}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		config.DataURL,
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

func PostHistory(data []receiverModel.WeatherDTO) error {

	if len(data) == 0 {
		return nil
	}

	client := &http.Client{Timeout: 10 * time.Second}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		config.BatchURL,
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
