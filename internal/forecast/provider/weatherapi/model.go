package weatherapi

type WeatherApiForecast struct {
	Location Location `json:"location"`
	Forecast Forecast `json:"forecast"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Forecast struct {
	ForecastDay []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date string `json:"date"`
	Hour []Hour `json:"hour"`
}

type Hour struct {
	Timestamp           int64   `json:"time_epoch"`
	Temperature         float64 `json:"temp_c"`
	Humidity            float64 `json:"humidity"`
	Precipitation       float64 `json:"precip_mm"`
	ApparentTemperature float64 `json:"feelslike_c"`
}
