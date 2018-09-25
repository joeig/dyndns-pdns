package main

import (
	"fmt"
	"github.com/joeig/go-powerdns"
	"log"
)

func addIPv4ResourceRecord(domain string, hostname string, ipv4 string, ttl int) error {
	return addPowerDNSResourceRecord(domain, hostname, "A", ipv4, ttl)
}

func addIPv6ResourceRecord(domain string, hostname string, ipv6 string, ttl int) error {
	return addPowerDNSResourceRecord(domain, hostname, "AAAA", ipv6, ttl)
}

func addPowerDNSResourceRecord(domain string, hostname string, recordType string, content string, ttl int) error {
	log.Printf("Calling PowerDNS (add) with domain='%s' hostname='%s' recordType='%s' content='%s' ttl=%d", domain, hostname, recordType, content, ttl)
	if Dry {
		log.Print("Dry run is enabled: Skipping calls to PowerDNS")
		return nil
	}
	pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, C.PowerDNS.APIKey)
	zone, err := pdns.GetZone(domain)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	name := fmt.Sprintf("%s.%s", hostname, domain)
	if err := zone.AddRecord(name, recordType, ttl, []string{content}); err != nil {
		log.Printf("Error changing %s record: %+v", recordType, err)
		return err
	}
	log.Print("Successfully created resource record")
	return nil
}

func deleteIPv4ResourceRecord(domain string, hostname string) error {
	return deletePowerDNSResourceRecord(domain, hostname, "A")
}

func deleteIPv6ResourceRecord(domain string, hostname string) error {
	return deletePowerDNSResourceRecord(domain, hostname, "AAAA")
}

func deletePowerDNSResourceRecord(domain string, hostname string, recordType string) error {
	log.Printf("Calling PowerDNS (delete) with domain='%s' hostname='%s' recordType='%s'", domain, hostname, recordType)
	if Dry {
		log.Print("Dry run is enabled: Skipping calls to PowerDNS")
		return nil
	}
	pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, C.PowerDNS.APIKey)
	zone, err := pdns.GetZone(domain)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	name := fmt.Sprintf("%s.%s", hostname, domain)
	if err := zone.DeleteRecord(name, recordType); err != nil {
		log.Printf("Error deleting %s record: %+v", recordType, err)
		return err
	}
	log.Print("Successfully deleted resource record")
	return nil
}
