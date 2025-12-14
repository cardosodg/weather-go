package config

import (
	"log"
	"os"
)

type ConfigDB struct {
	Bucket string
	Org    string
	Token  string
	Url    string
}

func LoadConfigDB() ConfigDB {
	required := []string{"INFLUXDB_ORG", "INFLUXDB_BUCKET", "INFLUXDB_ADMIN_TOKEN", "INFLUXDB_URL"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Fatal("Value not defined: ", v)
		}
	}

	return ConfigDB{
		Bucket: os.Getenv("INFLUXDB_BUCKET"),
		Org:    os.Getenv("INFLUXDB_ORG"),
		Token:  os.Getenv("INFLUXDB_ADMIN_TOKEN"),
		Url:    os.Getenv("INFLUXDB_URL"),
	}
}
