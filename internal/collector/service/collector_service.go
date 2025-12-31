package service

import (
	"WeatherTrack/internal/collector/config"
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

	url := fmt.Sprintf(config.OpenMeteoBaseURL, latitude, longitude, config.OpenMeteoParams)

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

	url := fmt.Sprintf(config.OpenMeteoHistURL, latitude, longitude, config.OpenMeteoParams)

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
