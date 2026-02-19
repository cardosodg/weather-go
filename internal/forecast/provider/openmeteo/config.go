package openmeteo

const (
	OpenMeteoTimeLayout = "2006-01-02T15:04"

	OpenMeteoParams      = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoForecastURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&forecast_days=2&hourly=%s&timezone=UTC"
)
