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
	var batch []WeatherData.WeatherData

	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"received_count": len(batch),
		"received":       batch,
	})
}
