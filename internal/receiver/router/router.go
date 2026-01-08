package router

import (
	IncomingData "WeatherTrack/internal/collector/model"
	"WeatherTrack/internal/receiver/config"
	"WeatherTrack/internal/receiver/database"
	WeatherData "WeatherTrack/internal/receiver/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *database.InfluxDB) {
	r.GET("/ping", healthHandler)
	r.POST("/data", weatherMeasureHandler(db))
	r.POST("/batch", weatherBatchHandler(db))
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func weatherMeasureHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var incomingData IncomingData.WeatherApiData

		if err := c.ShouldBindJSON(&incomingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ts, _ := time.Parse(time.RFC3339, incomingData.Current.Timestamp)

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

		if err := db.WriteBatch(batch); err != nil {
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
