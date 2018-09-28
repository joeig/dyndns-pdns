package pdns

import (
	"fmt"
	"github.com/joeig/go-powerdns"
	"log"
)

// DeleteIPv4ResourceRecord deletes an IPv4 resource record
func (p *PowerDNS) DeleteIPv4ResourceRecord(hostname string) error {
	return p.deletePowerDNSResourceRecord(hostname, "A")
}

// DeleteIPv6ResourceRecord deletes an IPv6 resource record
func (p *PowerDNS) DeleteIPv6ResourceRecord(hostname string) error {
	return p.deletePowerDNSResourceRecord(hostname, "AAAA")
}

func (p *PowerDNS) deletePowerDNSResourceRecord(hostname string, recordType string) error {
	log.Printf("Calling PowerDNS (delete) with domain='%s' hostname='%s' recordType='%s'", p.Zone, hostname, recordType)
	if p.Dry {
		log.Print("Dry run is enabled: Skipping calls to PowerDNS")
		return nil
	}
	pdns := powerdns.NewClient(p.BaseURL, p.VHost, p.APIKey)
	zone, err := pdns.GetZone(p.Zone)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	name := fmt.Sprintf("%s.%s", hostname, p.Zone)
	if err := zone.DeleteRecord(name, recordType); err != nil {
		log.Printf("Error deleting %s record: %+v", recordType, err)
		return err
	}
	log.Print("Successfully deleted resource record")
	return nil
}
