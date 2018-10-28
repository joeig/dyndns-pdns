package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// BuildVersion is set at linking time
var BuildVersion string

// BuildGitCommit is set at linking time
var BuildGitCommit string

// Dry prohibits calling any backend services
var Dry = false

// Debug enables verbose log output
var Debug = false

func main() {
	// Command line flags
	configFile := flag.String("config", "config.yml", "Configuration file")
	dryFlag := flag.Bool("dry", false, "Dry run (do not call any backend services)")
	debugFlag := flag.Bool("debug", false, "Debug mode")
	version := flag.Bool("version", false, "Prints the version name")
	flag.Parse()

	// Version
	if *version {
		fmt.Printf("Build Version: %s\n", BuildVersion)
		fmt.Printf("Build Git Commit: %s\n", BuildGitCommit)
		os.Exit(0)
	}

	// Initialize configuration
	parseConfig(&C, configFile)
	setDNSProvider(&dnsProvider)
	Dry = *dryFlag
	if Dry {
		log.Print("Dry run enabled")
	}
	C.PowerDNS.Dry = Dry
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
