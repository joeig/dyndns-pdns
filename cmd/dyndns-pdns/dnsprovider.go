package main

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
)

// DNSProviderPowerDNS sets the DNS provider to PowerDNS.
//
// This setting uses a PowerDNS backend.
const DNSProviderPowerDNS dnsprovider.Type = "powerDNS"

var activeDNSProvider dnsprovider.DNSProvider

func setDNSProvider(d *dnsprovider.DNSProvider) {
	switch C.DNSProviderType {
	case DNSProviderPowerDNS:
		*d = &C.PowerDNS
	default:
		panic(fmt.Errorf("invalid dnsProviderType \"%s\"", C.DNSProviderType))
	}
}
