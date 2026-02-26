package weatherapi

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetWeatherApiForecast(
	lat string,
	lon string,
	label string,
) ([]receiverModel.WeatherDTO, error) {

	var now = time.Now().UTC().Truncate(time.Hour)
	var apiResponse WeatherApiForecast
	var dto []receiverModel.WeatherDTO
	apiKey := loadKey()

	url := fmt.Sprintf(
		WeatherApiForecastURL,
		label,
		ForecastDays,
		apiKey,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return dto, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dto, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return dto, err
	}

	for _, day := range apiResponse.Forecast.ForecastDay {
		for _, hour := range day.Hour {
			ts := time.Unix(hour.Timestamp, 0).UTC()

			if ts.Before(now) {
				continue
			}

			aux := receiverModel.WeatherDTO{
				Location:      label,
				Timestamp:     ts,
				Temperature:   hour.Temperature,
				Humidity:      hour.Humidity,
				Precipitation: hour.Precipitation,
				ApparentTemp:  hour.ApparentTemperature,
			}

			dto = append(dto, aux)
		}
	}

	return dto, err
}
