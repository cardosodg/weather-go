package service

import (
	"WeatherTrack/internal/forecast/config"
	"WeatherTrack/internal/forecast/provider/openmeteo"
	"WeatherTrack/internal/forecast/provider/weatherapi"
	receiverModel "WeatherTrack/internal/receiver/model"
	"fmt"
)

func GetForecastWeather(
	latitude string,
	longitude string,
	locationName string,
) ([]receiverModel.WeatherDTO, error) {

	switch config.Provider {
	case "openmeteo":
		return openmeteo.GetOpenMeteoForecast(latitude, longitude, locationName)

	case "weatherapi":
		return weatherapi.GetWeatherApiForecast(latitude, longitude, locationName)

	default:
		return []receiverModel.WeatherDTO{},
			fmt.Errorf("unknown provider: %s", config.Provider)
	}
}
