package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

// Adds an unique request ID to every single Gin request
func requestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid4, err := uuid.NewRandom()
		if err != nil {
			log.Fatal("Unable to generate request Id")
			return
		}
		requestId := uuid4.String()
		c.Set("RequestId", requestId)
		c.Header("X-Request-Id", requestId)
		log.SetPrefix(fmt.Sprintf("[%s] ", requestId))
		log.Printf("Set request Id to \"%s\"", requestId)
		c.Next()
	}
}
