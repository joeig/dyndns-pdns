package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthStatus contains information regarding the healthiness of the application
type HealthStatus struct {
	ApplicationRunning bool `json:"applicationRunning"`
}

// Health Gin route
func Health(ctx *gin.Context) {
	hs := &HealthStatus{
		ApplicationRunning: true,
	}

	ctx.JSON(http.StatusOK, hs)
}
