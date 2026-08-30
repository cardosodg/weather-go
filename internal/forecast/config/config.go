package config

import "time"

const (
	// possible values: openmeteo or weatherapi
	Provider = "openmeteo"

	LocationsFile = "./locations.json"

	BaseReceiverURL = "http://receiver:8123"
	ForecastPath    = "/forecast"
	HealthPath      = "/ping"

	ForecastURL = BaseReceiverURL + ForecastPath
	HealthURL   = BaseReceiverURL + HealthPath

	ReceiverMaxRetries    = 10
	ReceiverRetryInterval = 10 * time.Second
)
