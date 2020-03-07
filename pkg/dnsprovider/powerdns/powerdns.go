package powerdns

import "github.com/joeig/go-powerdns/v2"

// PowerDNS defines the configuration for the PowerDNS authoritative server.
type PowerDNS struct {
	// BaseURL contains the URL to the PowerDNS API.
	//
	// Example: http://127.0.0.1:8080/
	BaseURL string `mapstructure:"baseURL"`

	// VHost contains the PowerDNS virtual host.
	//
	// Example: localhost
	VHost string `mapstructure:"vhost"`

	// APIKey contains the API key.
	APIKey string `mapstructure:"apiKey"`

	// Zone contains the DNS zone, which contains all maintained resource records.
	//
	// Example: dyn.example.com
	Zone string `mapstructure:"zone"`

	// DefaultTTL sets a default TTL value in seconds.
	//
	// Example: 10
	DefaultTTL uint32 `mapstructure:"defaultTTL"`

	// Dry toggles the simulation mode off and on.
	//
	// If this is enabled, no actual API calls are made.
	Dry bool `mapstructure:"dry"`
}

func (p *PowerDNS) setupPowerDNSClient() *powerdns.Client {
	headers := map[string]string{"X-API-Key": p.APIKey}
	return powerdns.NewClient(p.BaseURL, p.VHost, headers, nil)
}
