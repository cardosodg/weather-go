package weatherapi

import "os"

const (
	WeatherApiBaseURL = "https://api.weatherapi.com/v1"

	WeatherApiCurrentURL = WeatherApiBaseURL + "/current.json?q=%s&key=%s"
	WeatherApiHistURL    = WeatherApiBaseURL + "/history.json?q=%s&dt=%s&end_dt=%s&key=%s"
)

func loadKey() string {
	return os.Getenv("WEATHERAPIKEY")
}
