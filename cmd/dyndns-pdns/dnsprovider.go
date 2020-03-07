package main

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"github.com/joeig/dyndns-pdns/pkg/dnsprovider"
)

// DNSProviderPowerDNS sets the DNS provider to PowerDNS.
//
// This setting uses a PowerDNS backend.
const DNSProviderPowerDNS dnsprovider.Type = "powerDNS"

var activeDNSProvider dnsprovider.DNSProvider

func setDNSProvider(d *dnsprovider.DNSProvider) {
	switch yamlconfig.C.DNSProviderType {
	case DNSProviderPowerDNS:
		*d = &yamlconfig.C.PowerDNS
	default:
		panic(fmt.Errorf("invalid dnsProviderType \"%s\"", yamlconfig.C.DNSProviderType))
	}
}
