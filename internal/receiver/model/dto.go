package model

import "time"

type WeatherDTO struct {
	Location      string    `json:"location"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	ApparentTemp  float64   `json:"apparent_temp"`
	Timestamp     time.Time `json:"timestamp"`
	Precipitation float64   `json:"precipitation"`
}
