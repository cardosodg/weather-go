package router

import (
	IncomingData "WeatherTrack/internal/collector/model"
	WeatherData "WeatherTrack/internal/receiver/model"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", healthHandler)
	r.POST("/data", weatherMeasureHandler)
	r.POST("/batch", weatherBatchHandler)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func weatherMeasureHandler(c *gin.Context) {
	var incomingData IncomingData.WeatherApiData

	if err := c.ShouldBindJSON(&incomingData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := WeatherData.WeatherData{
		Location:     incomingData.Location,
		Timestamp:    incomingData.Current.Timestamp,
		Temperature:  incomingData.Current.Temperature,
		Humidity:     incomingData.Current.Humidity,
		Rain:         incomingData.Current.Rain,
		ApparentTemp: incomingData.Current.ApparentTemp,
	}

	if dbg, err := json.MarshalIndent(data, "", "  "); err == nil {
		log.Println(string(dbg))
	}
	c.JSON(http.StatusOK, gin.H{
		"received": data,
	})
}

func weatherBatchHandler(c *gin.Context) {
	var incomingHistory IncomingData.WeatherApiHistory
	var batch []WeatherData.WeatherData

	if err := c.ShouldBindJSON(&incomingHistory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch = make([]WeatherData.WeatherData, 0, len(incomingHistory.DataList.Timestamp))

	for i := range incomingHistory.DataList.Timestamp {
		wd := WeatherData.WeatherData{
			Location:     incomingHistory.Location,
			Temperature:  incomingHistory.DataList.Temperature[i],
			Humidity:     incomingHistory.DataList.Humidity[i],
			ApparentTemp: incomingHistory.DataList.ApparentTemp[i],
			Rain:         incomingHistory.DataList.Rain[i],
			Timestamp:    incomingHistory.DataList.Timestamp[i],
		}
		batch = append(batch, wd)

	}
	log.Println(batch)

	c.JSON(http.StatusOK, gin.H{
		"items":  len(batch),
		"status": "ok",
	})
}
