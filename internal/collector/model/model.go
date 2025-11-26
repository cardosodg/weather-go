package model

type WeatherApiData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
	Current   Current `json:"current"`
}

type Current struct {
	Timestamp    string  `json:"time"`
	Temperature  float64 `json:"temperature_2m"`
	Humidity     float64 `json:"relative_humidity_2m"`
	Rain         float64 `json:"rain"`
	ApparentTemp float64 `json:"apparent_temperature"`
}

type WeatherApiHistory struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Location  string   `json:"location"`
	DataList  DataList `json:"hourly"`
}

type DataList struct {
	Timestamp    []string  `json:"time"`
	Temperature  []float64 `json:"temperature_2m"`
	Humidity     []float64 `json:"relative_humidity_2m"`
	Rain         []float64 `json:"rain"`
	ApparentTemp []float64 `json:"apparent_temperature"`
}
