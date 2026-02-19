package service

import (
	"WeatherTrack/internal/collector/provider/openmeteo"
	receiverModel "WeatherTrack/internal/receiver/model"
)

func GetSingleWeather(
	latitude string,
	longitude string,
	locationName string,
) (receiverModel.WeatherDTO, error) {

	return openmeteo.GetOpenMeteoCurrent(
		latitude,
		longitude,
		locationName,
	)
}

func GetHistoryWeather(
	latitude string,
	longitude string,
	locationName string,
) ([]receiverModel.WeatherDTO, error) {

	return openmeteo.GetOpenMeteoHistory(
		latitude,
		longitude,
		locationName,
	)
}

// Future implementation
// func GetSingleWeather(
// 	latitude string,
// 	longitude string,
// 	locationName string,
// ) (receiverModel.WeatherDTO, error) {

// 	switch config.Provider {
// 	case "openmeteo":
// 		return provider.GetOpenMeteoCurrent(latitude, longitude, locationName)

// 	case "weatherapi":
// 		return provider.GetWeatherAPICurrent(latitude, longitude, locationName)

// 	default:
// 		return receiverModel.WeatherDTO{},
// 			fmt.Errorf("unknown provider: %s", config.Provider)
// 	}
// }
