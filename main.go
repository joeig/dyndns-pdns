package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var dry = false
var debug = false

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

func main() {
	// Command line flags
	configFile := flag.String("config", "config.yml", "Configuration file")
	dryFlag := flag.Bool("dry", false, "Dry run (do not call any backend services)")
	debugFlag := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	// Initialize configuration
	parseConfig(&C, configFile)
	dry = *dryFlag
	if dry {
		log.Print("Dry run enabled")
	}
	debug = *debugFlag
	if debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// Initialize Gin router
	router := getGinEngine()

	// Run server
	if C.Server.TLS.Enable {
		log.Fatal(http.ListenAndServeTLS(C.Server.ListenAddress, C.Server.TLS.CertFile, C.Server.TLS.KeyFile, router))
	}
	log.Fatal(http.ListenAndServe(C.Server.ListenAddress, router))
}
