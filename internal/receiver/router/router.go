package router

import (
	IncomingData "WeatherTrack/internal/collector/model"
	IncomingForecast "WeatherTrack/internal/forecast/model"
	"WeatherTrack/internal/receiver/config"
	"WeatherTrack/internal/receiver/database"
	WeatherData "WeatherTrack/internal/receiver/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *database.InfluxDB) {
	r.GET("/ping", healthHandler(db))
	r.POST("/data", weatherMeasureHandler(db))
	r.POST("/batch", weatherBatchHandler(db))
	r.POST("/forecast", weatherForecastHandler(db))
}

func healthHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"
		httpStatus := http.StatusOK

		if err := db.IsReady(); err != nil {
			status = "nok"
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"message":   "pong",
			"status":    status,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}

func weatherMeasureHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var incomingData IncomingData.WeatherApiData

		if err := c.ShouldBindJSON(&incomingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ts, _ := time.Parse(config.OpenMeteoTimeLayout, incomingData.Current.Timestamp)

		data := WeatherData.WeatherData{
			Location:     incomingData.Location,
			Timestamp:    ts,
			Temperature:  incomingData.Current.Temperature,
			Humidity:     incomingData.Current.Humidity,
			Rain:         incomingData.Current.Rain,
			ApparentTemp: incomingData.Current.ApparentTemp,
		}

		if err := db.WriteData(data); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	}
}

func weatherBatchHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var incomingHistory IncomingData.WeatherApiHistory
		var batch []WeatherData.WeatherData
		var now = time.Now().UTC()

		if err := c.ShouldBindJSON(&incomingHistory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		batch = make([]WeatherData.WeatherData, 0, len(incomingHistory.DataList.Timestamp))

		for i := range incomingHistory.DataList.Timestamp {
			ts, err := time.Parse(config.OpenMeteoTimeLayout, incomingHistory.DataList.Timestamp[i])
			if err != nil {
				log.Println(err)
				continue
			}

			if ts.After(now) {
				continue
			}

			wd := WeatherData.WeatherData{
				Location:     incomingHistory.Location,
				Temperature:  incomingHistory.DataList.Temperature[i],
				Humidity:     incomingHistory.DataList.Humidity[i],
				ApparentTemp: incomingHistory.DataList.ApparentTemp[i],
				Rain:         incomingHistory.DataList.Rain[i],
				Timestamp:    ts,
			}
			batch = append(batch, wd)

		}

		if err := db.WriteBatch(batch, "weather_readings"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":  len(batch),
			"status": "ok",
		})
	}
}

func weatherForecastHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var forecast IncomingForecast.WeatherApiForecast
		var batch []WeatherData.WeatherData
		var now = time.Now().UTC()

		if err := c.ShouldBindJSON(&forecast); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		batch = make([]WeatherData.WeatherData, 0, len(forecast.DataList.Timestamp))

		for i := range forecast.DataList.Timestamp {
			ts, err := time.Parse(config.OpenMeteoTimeLayout, forecast.DataList.Timestamp[i])
			if err != nil {
				log.Println(err)
				continue
			}

			if ts.Before(now) {
				continue
			}

			wd := WeatherData.WeatherData{
				Location:     forecast.Location,
				Temperature:  forecast.DataList.Temperature[i],
				Humidity:     forecast.DataList.Humidity[i],
				ApparentTemp: forecast.DataList.ApparentTemp[i],
				Rain:         forecast.DataList.Rain[i],
				Timestamp:    ts,
			}
			batch = append(batch, wd)

		}

		if err := db.WriteBatch(batch, "forecast_readings"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":  len(batch),
			"status": "ok",
		})
	}
}
