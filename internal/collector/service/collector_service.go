package service

import (
	"WeatherTrack/internal/collector/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetSingleWeather(
	latitude string,
	longitude string,
	localtionName string,
) (model.WeatherApiData, error) {

	var incoming model.WeatherApiData
	const params string = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	const baseURL string = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=%s"

	url := fmt.Sprintf(baseURL, latitude, longitude, params)

	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return incoming, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.WeatherApiData{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&incoming)

	incoming.Location = localtionName
	log.Println(incoming)

	return incoming, nil
}
