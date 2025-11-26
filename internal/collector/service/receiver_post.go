package service

import (
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
	log.Println(string(body))

	resp, err := client.Post(
		"http://localhost:8123/data",
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
	log.Println(string(body))

	resp, err := client.Post(
		"http://localhost:8123/batch",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
