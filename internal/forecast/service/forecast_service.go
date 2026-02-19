package service

import (
	"WeatherTrack/internal/forecast/provider/openmeteo"
	receiverModel "WeatherTrack/internal/receiver/model"
)

func GetForecastWeather(
	latitude string,
	longitude string,
	localtionName string,
) ([]receiverModel.WeatherDTO, error) {

	return openmeteo.GetOpenMeteoForecast(
		latitude,
		longitude,
		localtionName,
	)
}
