package router

import (
	"WeatherTrack/internal/receiver/model"
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

// https://api.open-meteo.com/v1/forecast?latitude=-20.3222&longitude=-40.3381&current=temperature_2m,relative_humidity_2m,rain,apparent_temperature
func weatherMeasureHandler(c *gin.Context) {
	var data model.WeatherData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"received": data,
	})
}

func weatherBatchHandler(c *gin.Context) {
	var batch []model.WeatherData

	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"received_count": len(batch),
		"received":       batch,
	})
}
