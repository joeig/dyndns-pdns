package main

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
)

// DNSProviderTypePowerDNS sets the DNS provider type to PowerDNS
const DNSProviderTypePowerDNS dnsprovider.DNSProviderType = "powerDNS"

var activeDNSProvider dnsprovider.DNSProvider

func setDNSProvider(d *dnsprovider.DNSProvider) {
	switch C.DNSProviderType {
	case DNSProviderTypePowerDNS:
		*d = &C.PowerDNS
	default:
		panic(fmt.Errorf("invalid dnsProviderType \"%s\"", C.DNSProviderType))
	}
}
