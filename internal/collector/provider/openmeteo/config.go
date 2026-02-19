package openmeteo

const (
	OpenMeteoTimeLayout = "2006-01-02T15:04"

	OpenMeteoParams  = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=%s&timezone=UTC"
	OpenMeteoHistURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&past_days=60&longitude=%s&minutely_15=%s&timezone=UTC"
)
