package service

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetSingleWeather(
	latitude string,
	longitude string,
	localtionName string,
) (model.WeatherApiForecast, error) {

	var incoming model.WeatherApiForecast

	url := fmt.Sprintf(config.OpenMeteoForecastURL, latitude, longitude, config.OpenMeteoParams)

	resp, err := http.Get(url)
	if err != nil {
		return incoming, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.WeatherApiForecast{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&incoming)

	incoming.Location = localtionName

	return incoming, nil
}
