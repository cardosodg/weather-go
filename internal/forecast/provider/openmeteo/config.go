package openmeteo

const (
	OpenMeteoTimeLayout = "2006-01-02T15:04"

	OpenMeteoParams      = "temperature_2m,relative_humidity_2m,rain,apparent_temperature"
	OpenMeteoForecastURL = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&forecast_days=2&hourly=%s&timezone=UTC"

	OpenMeteoGridSize = 10
)

var OpenMeteoWeights = []float64{
	0.40, // Center
	0.10, // North
	0.10, // South
	0.10, // East
	0.10, // West
	0.05, // Northeast
	0.05, // Northwest
	0.05, // Southeast
	0.05, // Southwest
}
