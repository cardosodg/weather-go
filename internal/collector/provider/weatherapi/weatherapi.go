package weatherapi

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetWeatherApiCurrent(
	lat string,
	lon string,
	label string,
) (receiverModel.WeatherDTO, error) {

	var apiResponse WeatherApiCurrent
	var dto receiverModel.WeatherDTO
	apiKey := loadKey()

	url := fmt.Sprintf(WeatherApiCurrentURL, label, apiKey)
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

	ts := time.Unix(apiResponse.Current.Timestamp, 0).UTC()

	dto = receiverModel.WeatherDTO{
		Location:      label,
		Timestamp:     ts,
		Temperature:   apiResponse.Current.Temperature,
		Humidity:      apiResponse.Current.Humidity,
		Precipitation: apiResponse.Current.Precipitation,
		ApparentTemp:  apiResponse.Current.ApparentTemperature,
	}

	return dto, nil
}

func GetWeatherApiHistory(
	lat string,
	lon string,
	label string,
) ([]receiverModel.WeatherDTO, error) {

	now := time.Now().UTC()
	oldHistDate := now.AddDate(0, 0, -7)

	var apiResponse WeatherApiHistory
	var dto []receiverModel.WeatherDTO
	apiKey := loadKey()

	url := fmt.Sprintf(
		WeatherApiHistURL,
		label,
		oldHistDate.Format("2006-01-02"),
		now.Format("2006-01-02"),
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

	for _, day := range apiResponse.History.HistoryDay {
		for _, hour := range day.Hour {
			ts := time.Unix(hour.Timestamp, 0).UTC()

			if ts.After(now) {
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
