package config

const (
	OpenMeteoParams  = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=%s"
	OpenMeteoHistURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&past_days=72&longitude=%s&hourly=%s"

	LocationsFile = "./locations.json"

	BaseReceiverURL = "http://receiver:8123"
	DataPath        = "/data"
	BatchPath       = "/batch"

	DataURL  = BaseReceiverURL + DataPath
	BatchURL = BaseReceiverURL + BatchPath
)
