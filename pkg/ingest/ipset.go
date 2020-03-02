package ingest

import "github.com/asaskevich/govalidator"

// IPSet is a combination of one IPv4 address and one IPv6 address
type IPSet struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

// HasIPv4 evaluates whether an IPSet contains a IPv4 address
func (ipSet *IPSet) HasIPv4() bool {
	return ipSet.IPv4 != ""
}

// IsIPv4 evaluates whether an IPSet contains a valid IPv4 address
func (ipSet *IPSet) IsIPv4() bool {
	return govalidator.IsIPv4(ipSet.IPv4)
}

// HasIPv6 evaluates whether an IPSet contains a IPv6 address
func (ipSet *IPSet) HasIPv6() bool {
	return ipSet.IPv6 != ""
}

// IsIPv6 evaluates whether an IPSet contains a valid IPv6 address
func (ipSet *IPSet) IsIPv6() bool {
	return govalidator.IsIPv6(ipSet.IPv6)
}
