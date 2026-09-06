package openmeteo

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"io"
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
		Precipitation       float64 `json:"precipitation"`
		ApparentTemperature float64 `json:"apparent_temperature"`
	} `json:"current"`
}

type openMeteoHistory struct {
	Minutely15 struct {
		Time                []string   `json:"time"`
		Temperature         []*float64 `json:"temperature_2m"`
		Humidity            []*float64 `json:"relative_humidity_2m"`
		Precipitation       []*float64 `json:"precipitation"`
		ApparentTemperature []*float64 `json:"apparent_temperature"`
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
	var totalTemp, totalHumidity, totalPrecip, totalApparent float64

	for i, resp := range responses {
		if i >= len(OpenMeteoWeights) {
			break
		}

		weight := OpenMeteoWeights[i]

		totalTemp += resp.Current.Temperature * weight
		totalHumidity += resp.Current.Humidity * weight
		totalPrecip += resp.Current.Precipitation * weight
		totalApparent += resp.Current.ApparentTemperature * weight
	}

	ts, _ := time.Parse(OpenMeteoTimeLayout, responses[0].Current.Time)

	dto = receiverModel.WeatherDTO{
		Location:      label,
		Timestamp:     ts,
		Temperature:   totalTemp,
		Humidity:      totalHumidity,
		Precipitation: totalPrecip,
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return dto, fmt.Errorf("openmeteo error [status %d]: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return dto, fmt.Errorf("failed to decode json: %v | raw response: %s", err, string(bodyBytes))
	}

	if len(apiResp) == 0 {
		return dto, fmt.Errorf("empty response from open-meteo current grid")
	}

	dto = EvaluateGridAverage(apiResp, label)

	return dto, nil
}

func evaluateHistoryGridAt(
	incoming []openMeteoHistory,
	index int,
	label string,
	ts time.Time,
) (receiverModel.WeatherDTO, bool) {
	var totalTemp, totalHumidity, totalPrecip, totalApparent float64

	for j, point := range incoming {
		if j >= len(OpenMeteoWeights) {
			break
		}

		if index >= len(point.Minutely15.Temperature) {
			return receiverModel.WeatherDTO{}, false
		}

		tempPtr := point.Minutely15.Temperature[index]
		humPtr := point.Minutely15.Humidity[index]
		precipPtr := point.Minutely15.Precipitation[index]
		appPtr := point.Minutely15.ApparentTemperature[index]

		if tempPtr == nil || humPtr == nil || precipPtr == nil || appPtr == nil {
			return receiverModel.WeatherDTO{}, false
		}

		weight := OpenMeteoWeights[j]

		totalTemp += *tempPtr * weight
		totalHumidity += *humPtr * weight
		totalPrecip += *precipPtr * weight
		totalApparent += *appPtr * weight
	}

	dto := receiverModel.WeatherDTO{
		Location:      label,
		Timestamp:     ts,
		Temperature:   totalTemp,
		Humidity:      totalHumidity,
		Precipitation: totalPrecip,
		ApparentTemp:  totalApparent,
	}

	return dto, true
}

func GetOpenMeteoHistory(
	lat string,
	lon string,
	label string,
) ([]receiverModel.WeatherDTO, error) {

	var apiResp []openMeteoHistory
	var result []receiverModel.WeatherDTO
	var now = time.Now().UTC()

	lats, lons, err := buildGrid(lat, lon)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf(
		OpenMeteoHistURL,
		lats,
		lons,
		OpenMeteoParams,
	)

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read history response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("openmeteo history error [status %d]: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return result, fmt.Errorf("failed to decode history json: %v | raw response: %s", err, string(bodyBytes))
	}

	if len(apiResp) == 0 {
		return result, fmt.Errorf("empty response from open-meteo history grid")
	}

	referenceTimestamps := apiResp[0].Minutely15.Time

	for i := range referenceTimestamps {
		ts, err := time.Parse(OpenMeteoTimeLayout, referenceTimestamps[i])
		if err != nil {
			continue
		}

		if ts.After(now) {
			continue
		}

		dto, ok := evaluateHistoryGridAt(apiResp, i, label, ts)
		if !ok {
			continue
		}

		result = append(result, dto)
	}

	return result, nil
}
