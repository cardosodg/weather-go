package openmeteo

import (
	receiverModel "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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

func GetOpenMeteoForecast(
	latitude string,
	longitude string,
	localtionName string,
) ([]receiverModel.WeatherDTO, error) {

	var incoming openMeteoForecast
	var forecast []receiverModel.WeatherDTO
	var now = time.Now().UTC()

	url := fmt.Sprintf(OpenMeteoForecastURL, latitude, longitude, OpenMeteoParams)

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

	incoming.Location = localtionName

	for i := range incoming.DataList.Timestamp {
		ts, err := time.Parse(OpenMeteoTimeLayout, incoming.DataList.Timestamp[i])
		if err != nil {
			log.Println(err)
			continue
		}

		if ts.Before(now) {
			continue
		}

		wd := receiverModel.WeatherDTO{
			Location:      incoming.Location,
			Temperature:   incoming.DataList.Temperature[i],
			Humidity:      incoming.DataList.Humidity[i],
			ApparentTemp:  incoming.DataList.ApparentTemp[i],
			Precipitation: incoming.DataList.Rain[i],
			Timestamp:     ts,
		}
		forecast = append(forecast, wd)

	}

	return forecast, nil
}
