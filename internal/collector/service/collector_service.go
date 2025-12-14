package service

import (
	"WeatherTrack/internal/collector/model"
	"encoding/json"
	"fmt"
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

	return incoming, nil
}

func GetHistoryWeather(
	latitude string,
	longitude string,
	localtionName string,
) (model.WeatherApiHistory, error) {
	var history model.WeatherApiHistory
	const params string = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	const baseURL string = "https://api.open-meteo.com/v1/forecast?latitude=%s&past_days=72&longitude=%s&hourly=%s"

	url := fmt.Sprintf(baseURL, latitude, longitude, params)

	resp, err := http.Get(url)
	if err != nil {
		return history, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return history, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&history)

	history.Location = localtionName

	return history, nil

}
