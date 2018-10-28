package powerdns

// PowerDNS defines the configuration for the PowerDNS authoritative server
type PowerDNS struct {
	BaseURL    string `mapstructure:"baseURL"`
	VHost      string `mapstructure:"vhost"`
	APIKey     string `mapstructure:"apiKey"`
	Zone       string `mapstructure:"zone"`
	DefaultTTL int    `mapstructure:"defaultTTL"`
	Dry        bool   `mapstructure:"dry"`
}
