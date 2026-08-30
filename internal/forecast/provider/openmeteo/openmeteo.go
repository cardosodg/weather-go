package openmeteo

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"log"
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

type openMeteoForecast struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
	DataList  struct {
		Timestamp    []string  `json:"time"`
		Temperature  []float64 `json:"temperature_2m"`
		Humidity     []float64 `json:"relative_humidity_2m"`
		Rain         []float64 `json:"rain"`
		ApparentTemp []float64 `json:"apparent_temperature"`
	} `json:"hourly"`
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

func evaluateForecastGridAt(
	incoming []openMeteoForecast,
	index int,
	locationName string,
	ts time.Time,
) (receiverModel.WeatherDTO, bool) {
	var totalTemp, totalHumidity, totalApparent, totalRain float64

	for j, point := range incoming {
		if j >= len(OpenMeteoWeights) {
			break
		}

		if index >= len(point.DataList.Temperature) {
			return receiverModel.WeatherDTO{}, false
		}

		weight := OpenMeteoWeights[j]

		totalTemp += point.DataList.Temperature[index] * weight
		totalHumidity += point.DataList.Humidity[index] * weight
		totalApparent += point.DataList.ApparentTemp[index] * weight
		totalRain += point.DataList.Rain[index] * weight
	}

	wd := receiverModel.WeatherDTO{
		Location:      locationName,
		Temperature:   totalTemp,
		Humidity:      totalHumidity,
		ApparentTemp:  totalApparent,
		Precipitation: totalRain,
		Timestamp:     ts,
	}

	return wd, true
}

func GetOpenMeteoForecast(
	latitude string,
	longitude string,
	locationName string,
) ([]receiverModel.WeatherDTO, error) {

	var incoming []openMeteoForecast
	var forecast []receiverModel.WeatherDTO
	var now = time.Now().UTC()

	lats, lons, err := buildGrid(latitude, longitude)
	if err != nil {
		return forecast, err
	}

	url := fmt.Sprintf(OpenMeteoForecastURL, lats, lons, OpenMeteoParams)

	resp, err := http.Get(url)
	if err != nil {
		return forecast, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return forecast, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&incoming); err != nil {
		return forecast, err
	}

	if len(incoming) == 0 {
		return forecast, fmt.Errorf("empty response from open-meteo grid")
	}

	referenceTimestamps := incoming[0].DataList.Timestamp

	for i := range referenceTimestamps {
		ts, err := time.Parse(OpenMeteoTimeLayout, referenceTimestamps[i])
		if err != nil {
			log.Println(err)
			continue
		}

		if ts.Before(now) {
			continue
		}

		wd, ok := evaluateForecastGridAt(incoming, i, locationName, ts)
		if !ok {
			continue
		}

		forecast = append(forecast, wd)
	}
	return forecast, nil
}
