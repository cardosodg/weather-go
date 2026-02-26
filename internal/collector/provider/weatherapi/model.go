package weatherapi

type WeatherApiCurrent struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type WeatherApiHistory struct {
	Location Location `json:"location"`
	History  History  `json:"forecast"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Current struct {
	Timestamp           int64   `json:"last_updated_epoch"`
	Temperature         float64 `json:"temp_c"`
	Humidity            float64 `json:"humidity"`
	Precipitation       float64 `json:"precip_mm"`
	ApparentTemperature float64 `json:"feelslike_c"`
}

type History struct {
	HistoryDay []HistoryDay `json:"forecastday"`
}

type HistoryDay struct {
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
