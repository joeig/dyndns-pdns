package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/ginmiddleware"
)

// Initializes the Gin engine
func setupGinEngine() *gin.Engine {
	router := gin.Default()
	router.Use(ginmiddleware.RequestIDMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/health", Health)
		v1.GET("/host/:name/sync", HostSync)
	}

	return router
}
