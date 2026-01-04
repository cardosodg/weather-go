package service

import (
	"WeatherTrack/internal/collector/config"
	"WeatherTrack/internal/collector/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func PostData(data model.WeatherApiData) error {
	client := &http.Client{Timeout: 5 * time.Second}
	body, _ := json.Marshal(data)
	log.Printf("Posting data: %s", string(body))

	resp, err := client.Post(
		config.DataURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func PostHistory(data model.WeatherApiHistory) error {
	client := &http.Client{Timeout: 5 * time.Second}
	body, _ := json.Marshal(data)
	log.Printf("Posting data: %s, size: %d", data.Location, len(data.DataList.Timestamp))

	resp, err := client.Post(
		config.BatchURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
