package main

import "fmt"

// DNSProviderType sets the DNS provider type
type DNSProviderType string

// DNSProviderTypePowerDNS sets the DNS provider type to PowerDNS
const DNSProviderTypePowerDNS DNSProviderType = "powerDNS"

// DNSProvider is an interface for basic DNS operations
type DNSProvider interface {
	AddIPv4ResourceRecord(hostname string, ipv4 string, ttl int) error
	AddIPv6ResourceRecord(hostname string, ipv6 string, ttl int) error
	DeleteIPv4ResourceRecord(hostname string) error
	DeleteIPv6ResourceRecord(hostname string) error
}

var dnsProvider DNSProvider

func setDNSProvider(d *DNSProvider) {
	switch C.DNSProviderType {
	case DNSProviderTypePowerDNS:
		*d = &C.PowerDNS
	default:
		panic(fmt.Errorf("invalid dnsProviderType \"%s\"", C.DNSProviderType))
	}
}
