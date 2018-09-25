package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config contains the primary configuration structure of the application
type Config struct {
	Server struct {
		ListenAddress string `mapstructure:"listenaddress"`
		TLS           struct {
			Enable   bool   `mapstructure:"enable"`
			CertFile string `mapstructure:"certFile"`
			KeyFile  string `mapstructure:"keyFile"`
		} `mapstructure:"tls"`
	} `mapstructure:"server"`
	PowerDNS struct {
		BaseURL string `mapstructure:"baseURL"`
		VHost   string `mapstructure:"vhost"`
		APIKey  string `mapstructure:"apiKey"`
	} `mapstructure:"powerDNS"`
	Miscellaneous struct {
		Zone       string `mapstructure:"zone"`
		DefaultTTL int    `mapstructure:"defaultTTL"`
	} `mapstructure:"miscellaneous"`
	KeyTable []Key `mapstructure:"keyTable"`
}

// C initializes the primary configuration of the application
var C Config

func parseConfig(config *Config, configFile *string) {
	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("%s", err))
	}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("%s", err))
	}
}

// IngestMode sets the IP address ingest mode
type IngestMode string

const (
	// IngestModeGetParameter sets the ingest mode to GET
	IngestModeGetParameter IngestMode = "getParameter"
	// IngestModeRemoteAddress sets the ingest mode to the remote address
	IngestModeRemoteAddress IngestMode = "remoteAddress"
)

// Key defines the structure of a certain key item
type Key struct {
	Name       string     `mapstructure:"name"`
	Enable     bool       `mapstructure:"enable"`
	Key        string     `mapstructure:"key"`
	HostName   string     `mapstructure:"hostName"`
	IngestMode IngestMode `mapstructure:"ingestMode"`
	TTL        int        `mapstructure:"ttl"`
}
