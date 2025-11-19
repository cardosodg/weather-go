package model

type Location struct {
	Name      string
	tag       string
	latitude  float64
	longitude float64
}

type WeatherData struct {
	Temperature  float64
	Humidity     float64
	ApparentTemp float64
	Timestamp    string
	Rain         float64
}
