package powerdns

import "github.com/joeig/go-powerdns/v2"

// PowerDNS defines the configuration for the PowerDNS authoritative server
type PowerDNS struct {
	BaseURL    string `mapstructure:"baseURL"`
	VHost      string `mapstructure:"vhost"`
	APIKey     string `mapstructure:"apiKey"`
	Zone       string `mapstructure:"zone"`
	DefaultTTL uint32 `mapstructure:"defaultTTL"`
	Dry        bool   `mapstructure:"dry"`
}

func (p *PowerDNS) setupPowerDNSClient() *powerdns.Client {
	headers := map[string]string{"X-API-Key": p.APIKey}
	return powerdns.NewClient(p.BaseURL, p.VHost, headers, nil)
}
