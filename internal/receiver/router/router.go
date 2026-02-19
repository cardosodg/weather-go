package router

import (
	"WeatherTrack/internal/receiver/database"
	"WeatherTrack/internal/receiver/model"
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
		var data model.WeatherDTO

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
		var batch []model.WeatherDTO

		if err := c.ShouldBindJSON(&batch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"items":   0,
				"message": err.Error(),
			})

			return
		}

		if len(batch) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"items":   0,
				"message": "empty batch",
			})

			return
		}

		if err := db.WriteBatch(batch, "weather_readings"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":   len(batch),
			"status":  "ok",
			"message": "ok",
		})
	}
}

func weatherForecastHandler(db *database.InfluxDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var batch []model.WeatherDTO

		if err := c.ShouldBindJSON(&batch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"items":   0,
				"message": err.Error(),
			})

			return
		}

		if len(batch) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"items":   0,
				"message": "empty batch",
			})

			return
		}

		if err := db.WriteBatch(batch, "forecast_readings"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":   len(batch),
			"status":  "ok",
			"message": "ok",
		})
	}
}
