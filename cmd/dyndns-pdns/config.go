package main

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider/powerdns"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"github.com/spf13/viper"
)

const (
	// IngestModeGetParameter sets the ingest mode to "GET parameter".
	//
	// The IP addresses can be updated by providing the new values within the query string.
	// Endpoint: /v1/host/<device name>/sync?key=<key>&ipv4=<IPv4 address>&ipv6=<IPv6 address>
	IngestModeGetParameter ingest.ModeType = "getParameter"

	// IngestModeRemoteAddress sets the ingest mode to "remote address".
	//
	// Instead of using values from the query string, this mode uses the remote address which was provided by the operating system's TCP/IP stack.
	// This might cause issues if you're running dyndns-pdns behind a reverse proxy.
	// Endpoint: /v1/host/<device name>/sync?key=<key>
	IngestModeRemoteAddress ingest.ModeType = "remoteAddress"
)

// Config contains the primary configuration structure of the application.
type Config struct {
	// Server contains the server configuration structure.
	Server Server `mapstructure:"server"`

	// Type sets the particular DNS provider type.
	//
	// Example: powerDNS
	DNSProviderType dnsprovider.Type `mapstructure:"dnsProviderType"`

	// PowerDNS sets the particular PowerDNS configuration, if the Type was set to "powerDNS".
	PowerDNS powerdns.PowerDNS `mapstructure:"powerDNS"`

	// KeyTable sets a list of host configurations.
	KeyTable []Key `mapstructure:"keyTable"`
}

// Server defines the structure of the server configuration.
type Server struct {
	// ListenAddress defines the local listener of the dyndns-pdns server.
	//
	// Example: 127.0.0.1:8000
	ListenAddress string `mapstructure:"listenaddress"`

	// TLS defines the particular TLS configuration of the dyndns-pdns server.
	TLS TLS `mapstructure:"tls"`
}

// TLS defines the structure of the TLS configuration.
type TLS struct {
	// Enable toggles the TLS mode off and on.
	Enable bool `mapstructure:"enable"`

	// CertFile is an absolute or relative path to the certificate file.
	CertFile string `mapstructure:"certFile"`

	// KeyFile is an absolute or relative path to the key file.
	KeyFile string `mapstructure:"keyFile"`
}

// Key defines the structure of a certain key item.
type Key struct {
	// Name is a human-friendly name of the particular host.
	//
	// Example: homeRouter
	Name string `mapstructure:"name"`

	// Enable toggles updates for the host off and on.
	Enable bool `mapstructure:"enable"`

	// Key contains the password for the device, which is required in order to update IP addresses.
	Key string `mapstructure:"key"`

	// HostName contains the host part of the maintained resource record.
	//
	// The device will be reachable via <HostName>.<Zone>.
	// Example: home-router
	HostName string `mapstructure:"hostName"`

	// Mode specifies the ingest mode.
	//
	// Example: getParameter
	IngestMode ingest.ModeType `mapstructure:"ingestMode"`

	// CleanUpMode specifies the clean up mode.
	//
	// Example: any
	CleanUpMode CleanUpModeType `mapstructure:"cleanUpMode"`

	// TTL overrides the default TTL value in seconds.
	//
	// Example: 5
	TTL uint32 `mapstructure:"ttl"`
}

// CleanUpModeType defines the clean up mode type.
type CleanUpModeType string

const (
	// CleanUpModeAny removes both A and AAAA resource records for the particular key item, even if only one IP address type was requested.
	//
	// Example: If an IPv4 address is given, this cleans all existing A (IPv4) and AAAA (IPv6) resource records.
	CleanUpModeAny CleanUpModeType = "any"

	// CleanUpModeRequestBased removes only the existing resource record which corresponds to the requested IP address type.
	//
	// Example: If an IPv4 address is given, this cleans only the existing A (IPv4) resource record while keeping the corresponding AAAA (IPv6) rersource record untouched.
	CleanUpModeRequestBased CleanUpModeType = "requested"
)

// C initializes the primary configuration of the application.
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
