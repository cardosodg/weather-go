package openmeteo

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	Lat float64
	Lon float64
}

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

func buildGrid(latitude string, longitude string) (string, string, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return "", "", err
	}
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return "", "", err
	}

	latRad := lat * math.Pi / 180.0
	delta_lat := OpenMeteoGridSize / 111.0
	delta_lon := OpenMeteoGridSize / (111.0 * math.Cos(latRad))

	points := []Point{
		{Lat: lat, Lon: lon},                         // Central
		{Lat: lat + delta_lat, Lon: lon},             // North
		{Lat: lat - delta_lat, Lon: lon},             // South
		{Lat: lat, Lon: lon + delta_lon},             // East
		{Lat: lat, Lon: lon - delta_lon},             // West
		{Lat: lat + delta_lat, Lon: lon + delta_lon}, // Northeast
		{Lat: lat + delta_lat, Lon: lon - delta_lon}, // Northwest
		{Lat: lat - delta_lat, Lon: lon + delta_lon}, // Southeast
		{Lat: lat - delta_lat, Lon: lon - delta_lon}, // Southwest
	}

	var latsBuilder, lonsBuilder strings.Builder

	for i, p := range points {
		if i > 0 {
			latsBuilder.WriteString(",")
			lonsBuilder.WriteString(",")
		}
		fmt.Fprintf(&latsBuilder, "%.6f", p.Lat)
		fmt.Fprintf(&lonsBuilder, "%.6f", p.Lon)
	}

	return latsBuilder.String(), lonsBuilder.String(), nil
}

func EvaluateGridAverage(responses []openMeteoCurrent, label string) receiverModel.WeatherDTO {
	var dto receiverModel.WeatherDTO
	var totalTemp, totalHumidity, totalRain, totalApparent float64

	for i, resp := range responses {
		if i >= len(OpenMeteoWeights) {
			break
		}

		weight := OpenMeteoWeights[i]

		totalTemp += resp.Current.Temperature * weight
		totalHumidity += resp.Current.Humidity * weight
		totalRain += resp.Current.Rain * weight
		totalApparent += resp.Current.ApparentTemperature * weight
	}

	ts, _ := time.Parse(OpenMeteoTimeLayout, responses[0].Current.Time)

	dto = receiverModel.WeatherDTO{
		Location:      label,
		Timestamp:     ts,
		Temperature:   totalTemp,
		Humidity:      totalHumidity,
		Precipitation: totalRain,
		ApparentTemp:  totalApparent,
	}

	return dto

}

func GetOpenMeteoCurrent(
	lat string,
	lon string,
	label string,
) (receiverModel.WeatherDTO, error) {

	var apiResp []openMeteoCurrent
	var dto receiverModel.WeatherDTO

	lats, lons, _ := buildGrid(lat, lon)

	url := fmt.Sprintf(
		OpenMeteoBaseURL,
		lats,
		lons,
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

	dto = EvaluateGridAverage(apiResp, label)

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
