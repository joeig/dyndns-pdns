package dnsprovider

// Type sets the DNS provider type
type Type string

// DNSProvider is an interface for basic DNS operations
type DNSProvider interface {
	AddIPv4ResourceRecord(hostname string, ipv4 string, ttl uint32) error
	AddIPv6ResourceRecord(hostname string, ipv6 string, ttl uint32) error
	DeleteIPv4ResourceRecord(hostname string) error
	DeleteIPv6ResourceRecord(hostname string) error
}
