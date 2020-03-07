package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"log"
	"os"
)

// BuildVersion is set at linking time
var BuildVersion string

// BuildGitCommit is set at linking time
var BuildGitCommit string

func printVersionAndExit() {
	fmt.Printf("Build Version: %s\n", BuildVersion)
	fmt.Printf("Build Git Commit: %s\n", BuildGitCommit)
	os.Exit(0)
}

// Dry prohibits calling any backend services
var Dry = false

func toggleDryMode() {
	if Dry {
		log.Print("Dry run enabled")
	}
	yamlconfig.C.PowerDNS.Dry = Dry
}

// Debug enables verbose log output
var Debug = false

func toggleDebugMode() {
	if Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}
}

func runServer(router *gin.Engine) {
	if yamlconfig.C.Server.TLS.Enable {
		log.Fatal(router.RunTLS(yamlconfig.C.Server.ListenAddress, yamlconfig.C.Server.TLS.CertFile, yamlconfig.C.Server.TLS.KeyFile))
	}
	log.Fatal(router.Run(yamlconfig.C.Server.ListenAddress))
}

func main() {
	configFile := flag.String("config", "config.yml", "Configuration file")
	dryFlag := flag.Bool("dry", false, "Dry run (do not call any backend services)")
	debugFlag := flag.Bool("debug", false, "Debug mode")
	version := flag.Bool("version", false, "Prints the version name")
	flag.Parse()

	if *version {
		printVersionAndExit()
	}

	yamlconfig.ParseConfig(&yamlconfig.C, configFile)
	setDNSProvider(&activeDNSProvider)

	Dry = *dryFlag
	toggleDryMode()

	Debug = *debugFlag
	toggleDebugMode()

	router := setupGinEngine()
	runServer(router)
}
