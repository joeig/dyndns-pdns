package main

import "github.com/gin-gonic/gin"

// Initializes the Gin engine
func getGinEngine() *gin.Engine {
	router := gin.Default()
	router.Use(requestIdMiddleware())
	v1 := router.Group("/v1")
	{
		v1.GET("/health", Health)
		v1.GET("/host/:name/sync", HostSync)
	}
	return router
}
