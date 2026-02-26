package weatherapi

import "os"

const (
	ForecastDays          = 3
	WeatherApiBaseURL     = "https://api.weatherapi.com/v1"
	WeatherApiForecastURL = WeatherApiBaseURL + "/forecast.json?q=%s&days=%d&key=%s"
)

func loadKey() string {
	return os.Getenv("WEATHERAPIKEY")
}
