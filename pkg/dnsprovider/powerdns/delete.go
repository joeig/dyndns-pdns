package powerdns

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/joeig/go-powerdns/v2"
	"log"
)

// DeleteIPv4ResourceRecord deletes an IPv4 resource record
func (p *PowerDNS) DeleteIPv4ResourceRecord(hostname string) error {
	return p.deletePowerDNSResourceRecord(hostname, powerdns.RRTypeA)
}

// DeleteIPv6ResourceRecord deletes an IPv6 resource record
func (p *PowerDNS) DeleteIPv6ResourceRecord(hostname string) error {
	return p.deletePowerDNSResourceRecord(hostname, powerdns.RRTypeAAAA)
}

func (p *PowerDNS) deletePowerDNSResourceRecord(hostname string, recordType powerdns.RRType) error {
	if !govalidator.IsDNSName(hostname) {
		return &powerdns.Error{}
	}

	log.Printf("Calling PowerDNS (delete) with domain='%s' hostname='%s' recordType='%s'", p.Zone, hostname, recordType)

	if p.Dry {
		log.Print("Dry run is enabled: Skipping calls to PowerDNS")
		return nil
	}

	pdns := p.setupPowerDNSClient()

	name := fmt.Sprintf("%s.%s", hostname, p.Zone)
	log.Printf("Generated name='%s'", name)

	if err := pdns.Records.Delete(p.Zone, name, recordType); err != nil {
		log.Printf("Error deleting %s record: %+v", recordType, err)
		return err
	}

	log.Print("Successfully deleted resource record")
	return nil
}
