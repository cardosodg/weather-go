package service

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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
		log.Println("XXXXXXXXXXXXXXXX")
		log.Println(err)
		log.Println("XXXXXXXXXXXXXXXX")
		return err
	}

	defer resp.Body.Close()

	return nil
}
