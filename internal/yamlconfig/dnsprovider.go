package yamlconfig

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
)

// DNSProviderPowerDNS sets the DNS provider to PowerDNS.
//
// This setting uses a PowerDNS backend.
const DNSProviderPowerDNS dnsprovider.Type = "powerDNS"

// ActiveDNSProvider contains the currently activated DNS provider
var ActiveDNSProvider dnsprovider.DNSProvider

// SetDNSProvider configures the DNS provider set by the configuration
func SetDNSProvider(d *dnsprovider.DNSProvider) {
	switch C.DNSProviderType {
	case DNSProviderPowerDNS:
		*d = &C.PowerDNS
	default:
		panic(fmt.Errorf("invalid dnsProviderType \"%s\"", C.DNSProviderType))
	}
}
