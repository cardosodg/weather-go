package config

import "time"

const (
	OpenMeteoParams      = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoForecastURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&forecast_days=2&hourly=%s&timezone=UTC"

	LocationsFile = "./locations.json"

	BaseReceiverURL = "http://receiver:8123"
	ForecastPath    = "/forecast"
	HealthPath      = "/ping"

	ForecastURL = BaseReceiverURL + ForecastPath
	HealthURL   = BaseReceiverURL + HealthPath

	ReceiverMaxRetries    = 10
	ReceiverRetryInterval = 10 * time.Second
)
