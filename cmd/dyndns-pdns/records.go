package main

import (
	"github.com/joeig/dyndns-pdns/internal/ginresponse"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"log"
	"net/http"
)

func cleanUpOutdatedResourceRecords(ipSet *ingest.IPSet, keyItem *yamlconfig.Key, hostname string) error {
	if keyItem.CleanUpMode == yamlconfig.CleanUpModeAny || (keyItem.CleanUpMode == yamlconfig.CleanUpModeRequestBased && ipSet.HasIPv4()) {
		log.Print("Cleaning up any previously created IPv4 resource records")

		if err := yamlconfig.ActiveDNSProvider.DeleteIPv4ResourceRecord(hostname); err != nil {
			log.Printf("%+v", err)
			return &ginresponse.HTTPError{Message: "IPv4 record deletion failed", HTTPErrorCode: http.StatusInternalServerError}
		}
	} else {
		log.Print("Skipping clean up of previously created IPv4 resource records")
	}

	if keyItem.CleanUpMode == yamlconfig.CleanUpModeAny || (keyItem.CleanUpMode == yamlconfig.CleanUpModeRequestBased && ipSet.HasIPv6()) {
		log.Print("Cleaning up any previously created IPv6 resource records")

		if err := yamlconfig.ActiveDNSProvider.DeleteIPv6ResourceRecord(hostname); err != nil {
			log.Printf("%+v", err)
			return &ginresponse.HTTPError{Message: "IPv6 record deletion failed", HTTPErrorCode: http.StatusInternalServerError}
		}
	} else {
		log.Print("Skipping clean up of previously created IPv6 resource records")
	}

	return nil
}

func createNewIPv4ResourceRecord(ipSet *ingest.IPSet, keyItem *yamlconfig.Key, hostname string) error {
	log.Print("Creating IPv4 resource records")

	if err := yamlconfig.ActiveDNSProvider.AddIPv4ResourceRecord(hostname, ipSet.IPv4, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv4 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewIPv6ResourceRecord(ipSet *ingest.IPSet, keyItem *yamlconfig.Key, hostname string) error {
	log.Print("Creating IPv6 resource records")

	if err := yamlconfig.ActiveDNSProvider.AddIPv6ResourceRecord(hostname, ipSet.IPv6, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv6 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewResourceRecords(ipSet *ingest.IPSet, keyItem *yamlconfig.Key) (*HostSyncObject, error) {
	if ipSet.HasIPv4() {
		if err := createNewIPv4ResourceRecord(ipSet, keyItem, keyItem.HostName); err != nil {
			return nil, err
		}
	}

	if ipSet.HasIPv6() {
		if err := createNewIPv6ResourceRecord(ipSet, keyItem, keyItem.HostName); err != nil {
			return nil, err
		}
	}

	if ipSet.HasIPv4() || ipSet.HasIPv6() {
		hostSyncObject := &HostSyncObject{
			HostName:    keyItem.HostName,
			IngestMode:  keyItem.IngestMode,
			CleanUpMode: keyItem.CleanUpMode,
			TTL:         int(keyItem.TTL),
			IPv4:        ipSet.IPv4,
			IPv6:        ipSet.IPv6,
		}
		return hostSyncObject, nil
	}

	if len(keyItem.DynamicKeys) > 0 {
		return nil, nil
	}

	return nil, &ginresponse.HTTPError{Message: "No IP addresses to create resource records for", HTTPErrorCode: http.StatusBadRequest}
}
