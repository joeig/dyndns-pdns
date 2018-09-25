package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var Dry = false
var Debug = false

func main() {
	// Command line flags
	configFile := flag.String("config", "config.yml", "Configuration file")
	dryFlag := flag.Bool("dry", false, "Dry run (do not call any backend services)")
	debugFlag := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	// Initialize configuration
	parseConfig(&C, configFile)
	Dry = *dryFlag
	if Dry {
		log.Print("Dry run enabled")
	}
	Debug = *debugFlag
	if Debug {
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
