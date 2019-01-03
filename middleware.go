package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

// Adds an unique request ID to every single Gin request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid4, err := uuid.NewRandom()
		if err != nil {
			log.Fatal("Unable to generate request Id")
			return
		}
		requestID := uuid4.String()
		c.Set("RequestId", requestID)
		c.Header("X-Request-Id", requestID)
		log.SetPrefix(fmt.Sprintf("[%s] ", requestID))
		log.Printf("Set request Id to \"%s\"", requestID)
		c.Next()
	}
}
