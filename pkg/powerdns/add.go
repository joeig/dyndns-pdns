package powerdns

import (
	"fmt"
	"github.com/joeig/dyndns-pdns/pkg/tools"
	"github.com/joeig/go-powerdns/v2"
	"log"
)

// AddIPv4ResourceRecord adds a new IPv4 resource record
func (p *PowerDNS) AddIPv4ResourceRecord(hostname string, ipv4 string, ttl uint32) error {
	return p.addPowerDNSResourceRecord(hostname, powerdns.RRTypeA, ipv4, ttl)
}

// AddIPv6ResourceRecord adds a new IPv6 resource record
func (p *PowerDNS) AddIPv6ResourceRecord(hostname string, ipv6 string, ttl uint32) error {
	return p.addPowerDNSResourceRecord(hostname, powerdns.RRTypeAAAA, ipv6, ttl)
}

func (p *PowerDNS) addPowerDNSResourceRecord(hostname string, recordType powerdns.RRType, content string, ttl uint32) error {
	log.Printf("Calling PowerDNS (add) with domain='%s' hostname='%s' recordType='%s' content='%s' ttl=%d", p.Zone, hostname, recordType, content, ttl)

	if p.Dry {
		log.Print("Dry run is enabled: Skipping calls to PowerDNS")
		return nil
	}

	headers := map[string]string{"X-API-Key": p.APIKey}
	pdns := powerdns.NewClient(p.BaseURL, p.VHost, headers, nil)

	name := fmt.Sprintf("%s.%s", hostname, p.Zone)
	thisTTL := tools.GetTTL(ttl, p.DefaultTTL)
	log.Printf("Generated name='%s' ttl=%d", name, thisTTL)

	if err := pdns.Records.Add(p.Zone, name, recordType, thisTTL, []string{content}); err != nil {
		log.Printf("Error changing %s record: %+v", recordType, err)
		return err
	}

	log.Print("Successfully created resource record")
	return nil
}
