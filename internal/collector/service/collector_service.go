package service

import (
	"WeatherTrack/internal/collector/config"
	"WeatherTrack/internal/collector/provider/openmeteo"
	"WeatherTrack/internal/collector/provider/weatherapi"
	receiverModel "WeatherTrack/internal/receiver/model"
	"fmt"
)

func GetSingleWeather(
	latitude string,
	longitude string,
	locationName string,
) (receiverModel.WeatherDTO, error) {

	switch config.Provider {
	case "openmeteo":
		return openmeteo.GetOpenMeteoCurrent(latitude, longitude, locationName)

	case "weatherapi":
		return weatherapi.GetWeatherApiCurrent(latitude, longitude, locationName)

	default:
		return receiverModel.WeatherDTO{},
			fmt.Errorf("unknown provider: %s", config.Provider)
	}
}

func GetHistoryWeather(
	latitude string,
	longitude string,
	locationName string,
) ([]receiverModel.WeatherDTO, error) {

	switch config.Provider {
	case "openmeteo":
		return openmeteo.GetOpenMeteoHistory(latitude, longitude, locationName)

	case "weatherapi":
		return weatherapi.GetWeatherApiHistory(latitude, longitude, locationName)

	default:
		return []receiverModel.WeatherDTO{},
			fmt.Errorf("unknown provider: %s", config.Provider)
	}
}
