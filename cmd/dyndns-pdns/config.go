package main

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider/powerdns"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"github.com/spf13/viper"
)

const (
	// IngestModeGetParameter sets the ingest mode to GET
	IngestModeGetParameter ingest.IngestModeType = "getParameter"
	// IngestModeRemoteAddress sets the ingest mode to the remote address
	IngestModeRemoteAddress ingest.IngestModeType = "remoteAddress"
)

// Config contains the primary configuration structure of the application
type Config struct {
	Server          Server                      `mapstructure:"server"`
	DNSProviderType dnsprovider.DNSProviderType `mapstructure:"dnsProviderType"`
	PowerDNS        powerdns.PowerDNS           `mapstructure:"powerDNS"`
	KeyTable        []Key                       `mapstructure:"keyTable"`
}

// Server defines the structure of the server configuration
type Server struct {
	ListenAddress string `mapstructure:"listenaddress"`
	TLS           TLS    `mapstructure:"tls"`
}

// TLS defines the structure of the TLS configuration
type TLS struct {
	Enable   bool   `mapstructure:"enable"`
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
}

// Key defines the structure of a certain key item
type Key struct {
	Name       string                `mapstructure:"name"`
	Enable     bool                  `mapstructure:"enable"`
	Key        string                `mapstructure:"key"`
	HostName   string                `mapstructure:"hostName"`
	IngestMode ingest.IngestModeType `mapstructure:"ingestMode"`
	TTL        uint32                `mapstructure:"ttl"`
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
