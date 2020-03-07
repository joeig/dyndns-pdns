package main

import (
	"github.com/joeig/dyndns-pdns/internal/ginresponse"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"log"
	"net/http"
)

func cleanUpOutdatedResourceRecords(ipSet *ingest.IPSet, keyItem *yamlconfig.Key) error {
	if keyItem.CleanUpMode == yamlconfig.CleanUpModeAny || (keyItem.CleanUpMode == yamlconfig.CleanUpModeRequestBased && ipSet.HasIPv4()) {
		log.Print("Cleaning up any previously created IPv4 resource records")

		if err := activeDNSProvider.DeleteIPv4ResourceRecord(keyItem.HostName); err != nil {
			log.Printf("%+v", err)
			return &ginresponse.HTTPError{Message: "IPv4 record deletion failed", HTTPErrorCode: http.StatusInternalServerError}
		}
	} else {
		log.Print("Skipping clean up of previously created IPv4 resource records")
	}

	if keyItem.CleanUpMode == yamlconfig.CleanUpModeAny || (keyItem.CleanUpMode == yamlconfig.CleanUpModeRequestBased && ipSet.HasIPv6()) {
		log.Print("Cleaning up any previously created IPv6 resource records")

		if err := activeDNSProvider.DeleteIPv6ResourceRecord(keyItem.HostName); err != nil {
			log.Printf("%+v", err)
			return &ginresponse.HTTPError{Message: "IPv6 record deletion failed", HTTPErrorCode: http.StatusInternalServerError}
		}
	} else {
		log.Print("Skipping clean up of previously created IPv6 resource records")
	}

	return nil
}

func createNewIPv4ResourceRecord(ipSet *ingest.IPSet, keyItem *yamlconfig.Key) error {
	log.Print("Creating IPv4 resource records")

	if err := activeDNSProvider.AddIPv4ResourceRecord(keyItem.HostName, ipSet.IPv4, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv4 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewIPv6ResourceRecord(ipSet *ingest.IPSet, keyItem *yamlconfig.Key) error {
	log.Print("Creating IPv6 resource records")

	if err := activeDNSProvider.AddIPv6ResourceRecord(keyItem.HostName, ipSet.IPv6, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv6 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewResourceRecords(ipSet *ingest.IPSet, keyItem *yamlconfig.Key) error {
	if ipSet.HasIPv4() {
		if err := createNewIPv4ResourceRecord(ipSet, keyItem); err != nil {
			return err
		}
	}

	if ipSet.HasIPv6() {
		if err := createNewIPv6ResourceRecord(ipSet, keyItem); err != nil {
			return err
		}
	}

	return nil
}
