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

	if err := isReady(db); err != nil {
		return nil, err
	}

	return db, nil
}

func isReady(db *InfluxDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Client.Ready(ctx)

	return err
}

func parseTimestamp(ts string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	for _, layout := range formats {
		t, err := time.Parse(layout, ts)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("timestamp inválido: %s", ts)
}

func (db *InfluxDB) WriteData(data model.WeatherData) error {
	writeAPI := db.Client.WriteAPIBlocking(db.Org, db.Bucket)

	p := influxdb2.NewPoint(
		"weather_readings",
		map[string]string{
			"location": data.Location,
		},
		map[string]any{
			"temperature":   data.Temperature,
			"humidity":      data.Humidity,
			"rain":          data.Rain,
			"apparent_temp": data.ApparentTemp,
		},
		data.Timestamp,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return writeAPI.WritePoint(ctx, p)
}

func (db *InfluxDB) WriteBatch(data []model.WeatherData, measurement string) error {
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
				"rain":          item.Rain,
			},
			item.Timestamp,
		)
		points = append(points, p)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return writeAPI.WritePoint(ctx, points...)

}
