package openmeteo

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type openMeteoCurrent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time                string  `json:"time"`
		Temperature         float64 `json:"temperature_2m"`
		Humidity            float64 `json:"relative_humidity_2m"`
		Rain                float64 `json:"rain"`
		ApparentTemperature float64 `json:"apparent_temperature"`
	} `json:"current"`
}

type openMeteoHistory struct {
	Minutely15 struct {
		Time                []string  `json:"time"`
		Temperature         []float64 `json:"temperature_2m"`
		Humidity            []float64 `json:"relative_humidity_2m"`
		Rain                []float64 `json:"rain"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
	} `json:"minutely_15"`
}

func GetOpenMeteoCurrent(
	lat string,
	lon string,
	label string,
) (receiverModel.WeatherDTO, error) {

	var apiResp openMeteoCurrent
	var dto receiverModel.WeatherDTO

	url := fmt.Sprintf(
		OpenMeteoBaseURL,
		lat,
		lon,
		OpenMeteoParams,
	)

	resp, err := http.Get(url)
	if err != nil {
		return dto, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return dto, err
	}

	// 🔥 Parse timestamp aqui
	ts, err := time.Parse(OpenMeteoTimeLayout, apiResp.Current.Time)
	if err != nil {
		return dto, err
	}

	// 🔥 Mapping para DTO
	dto = receiverModel.WeatherDTO{
		Location:      label,
		Timestamp:     ts,
		Temperature:   apiResp.Current.Temperature,
		Humidity:      apiResp.Current.Humidity,
		Precipitation: apiResp.Current.Rain,
		ApparentTemp:  apiResp.Current.ApparentTemperature,
	}

	return dto, nil
}

func GetOpenMeteoHistory(
	lat string,
	lon string,
	label string,
) ([]receiverModel.WeatherDTO, error) {

	var apiResp openMeteoHistory
	var result []receiverModel.WeatherDTO
	var now = time.Now().UTC()

	url := fmt.Sprintf(
		OpenMeteoHistURL,
		lat,
		lon,
		OpenMeteoParams,
	)

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return result, err
	}

	for i := range apiResp.Minutely15.Time {

		ts, err := time.Parse(OpenMeteoTimeLayout, apiResp.Minutely15.Time[i])
		if err != nil {
			continue
		}

		if ts.After(now) {
			continue
		}

		dto := receiverModel.WeatherDTO{
			Location:      label,
			Timestamp:     ts,
			Temperature:   apiResp.Minutely15.Temperature[i],
			Humidity:      apiResp.Minutely15.Humidity[i],
			Precipitation: apiResp.Minutely15.Rain[i],
			ApparentTemp:  apiResp.Minutely15.ApparentTemperature[i],
		}

		result = append(result, dto)
	}

	return result, nil
}
