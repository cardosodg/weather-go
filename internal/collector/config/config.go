package config

import "time"

const (
	OpenMeteoParams  = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=%s&timezone=UTC"
	OpenMeteoHistURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&past_days=60&longitude=%s&minutely_15=%s&timezone=UTC"

	OpenMeteoWaitTime = 15

	LocationsFile = "./locations.json"

	BaseReceiverURL = "http://receiver:8123"
	DataPath        = "/data"
	BatchPath       = "/batch"
	HealthPath      = "/ping"

	DataURL   = BaseReceiverURL + DataPath
	BatchURL  = BaseReceiverURL + BatchPath
	HealthURL = BaseReceiverURL + HealthPath

	ReceiverMaxRetries    = 10
	ReceiverRetryInterval = 10 * time.Second
)
