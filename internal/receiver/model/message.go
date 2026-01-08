package model

import "time"

type WeatherData struct {
	Location     string    `json:"location"`
	Temperature  float64   `json:"temperature_2m"`
	Humidity     float64   `json:"relative_humidity_2m"`
	ApparentTemp float64   `json:"apparent_temperature"`
	Timestamp    time.Time `json:"time"`
	Rain         float64   `json:"rain"`
}
