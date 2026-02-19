package config

import "time"

const (
	// possible values: openmeteo or weatherapi
	Provider = "openmeteo"

	CollectorWaitTime = 15

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
