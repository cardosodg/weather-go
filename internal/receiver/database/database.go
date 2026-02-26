package database

import (
	"WeatherTrack/internal/receiver/config"
	"WeatherTrack/internal/receiver/model"
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDB struct {
	Client influxdb2.Client
	Org    string
	Bucket string
}

func Initialize() (*InfluxDB, error) {
	configDB := config.LoadConfigDB()
	bucket := configDB.Bucket
	org := configDB.Org
	token := configDB.Token
	url := configDB.Url

	client := influxdb2.NewClient(url, token)
	db := &InfluxDB{
		Client: client,
		Org:    org,
		Bucket: bucket,
	}

	err := db.IsReady()

	return db, err
}

func (db *InfluxDB) IsReady() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queryAPI := db.Client.QueryAPI(db.Org)

	result, err := queryAPI.Query(ctx, "buckets() |> limit(n:1)")
	if err != nil {
		return fmt.Errorf("influxdb not ready: %w", err)
	}

	if result.Err() != nil {
		return fmt.Errorf("flux error: %w", result.Err())
	}

	return nil
}

func (db *InfluxDB) WriteData(data model.WeatherDTO) error {
	writeAPI := db.Client.WriteAPIBlocking(db.Org, db.Bucket)

	p := influxdb2.NewPoint(
		"weather_readings",
		map[string]string{
			"location": data.Location,
		},
		map[string]any{
			"temperature":   data.Temperature,
			"humidity":      data.Humidity,
			"precipitation": data.Precipitation,
			"apparent_temp": data.ApparentTemp,
		},
		data.Timestamp,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return writeAPI.WritePoint(ctx, p)
}

func (db *InfluxDB) WriteBatch(data []model.WeatherDTO, measurement string) error {
	if len(data) == 0 {
		return nil
	}

	writeAPI := db.Client.WriteAPIBlocking(db.Org, db.Bucket)

	points := make([]*write.Point, 0, len(data))

	for _, item := range data {

		p := influxdb2.NewPoint(
			measurement,
			map[string]string{
				"location": item.Location,
			},
			map[string]any{
				"temperature":   item.Temperature,
				"humidity":      item.Humidity,
				"apparent_temp": item.ApparentTemp,
				"precipitation": item.Precipitation,
			},
			item.Timestamp,
		)
		points = append(points, p)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return writeAPI.WritePoint(ctx, points...)

}

func (db *InfluxDB) DeleteMeasurement(measurement string, location string) error {
	if measurement == "" || location == "" {
		return fmt.Errorf("measurement and location cannot be empty")
	}

	deleteAPI := db.Client.DeleteAPI()
	start := time.Unix(0, 0)
	stop := time.Now().AddDate(200, 0, 0)

	predicate := fmt.Sprintf(
		`_measurement="%s" AND location="%s"`,
		measurement,
		location,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := deleteAPI.DeleteWithName(ctx, db.Org, db.Bucket, start, stop, predicate)
	if err != nil {
		return fmt.Errorf("failed to delete measurement %s: %w", measurement, err)
	}

	return nil
}
