package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
	"net"
	"net/http"
)

// HostSyncPayload returns a payload containing HostSyncObjects
type HostSyncPayload struct {
	HostSyncObjects []*HostSyncObject `json:"hostSyncObjects"`
}

// HostSyncObject contains a payload for the requester in order to identify the values that have been stored
type HostSyncObject struct {
	HostName   string     `json:"hostName"`
	IngestMode IngestMode `json:"ingestMode"`
	TTL        int        `json:"ttl"`
	IPv4       string     `json:"ipv4"`
	IPv6       string     `json:"ipv6"`
}

// HostSync Gin route
func HostSync(c *gin.Context) {
	// Disable HTTP caches
	c.Header("Cache-Control", "no-cache")

	// Get "name"
	name, err := checkHost(c.Param("name"))
	log.Printf("Received name=\"%s\"", name)
	if err != nil {
		errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: err.Error()}}}
		log.Printf("%+v", errPayload)
		c.JSON(http.StatusUnauthorized, errPayload)
		return
	}

	// Get "key"
	key, err := checkKey(c.Query("key"))
	log.Printf("Received key=\"%s\"", key)
	if err != nil {
		errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: err.Error()}}}
		log.Printf("%+v", errPayload)
		c.JSON(http.StatusUnauthorized, errPayload)
		return
	}

	// Look up the corresponding key item
	keyItem, err := checkAuthorization(C.KeyTable, name, key)
	if err != nil {
		errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: err.Error()}}}
		log.Printf("%+v", errPayload)
		c.JSON(http.StatusForbidden, errPayload)
		return
	}
	log.Printf("Found key item: %+v", keyItem)

	// Get IP addresses
	var ipv4 string
	var ipv6 string
	switch keyItem.IngestMode {
	case IngestModeGetParameter:
		log.Printf("Processing ingest for %+v mode", IngestModeGetParameter)
		ipv4 = c.Query("ipv4")
		ipv6 = c.Query("ipv6")
		log.Printf("Received ipv4=\"%s\" ipv6=\"%s\"", ipv4, ipv6)
		if ipv4 == "" && ipv6 == "" {
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv4 as well as IPv6 parameter missing"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusBadRequest, payload)
			return
		}
		if ipv4 != "" && !govalidator.IsIPv4(ipv4) {
			ipv4 = ""
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv4 address invalid"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusBadRequest, payload)
			return
		}
		if ipv6 != "" && !govalidator.IsIPv6(ipv6) {
			ipv6 = ""
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv6 address invalid"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusBadRequest, payload)
			return
		}
		break
	case IngestModeRemoteAddress:
		log.Printf("Processing ingest for %+v mode", IngestModeRemoteAddress)
		// Under certain circumstances, RemoteAddr contains also the port number
		address, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			address = c.Request.RemoteAddr
		}
		log.Printf("Received address=\"%s\"", address)
		if govalidator.IsIPv4(address) {
			ipv4 = address
			ipv6 = ""
		} else if govalidator.IsIPv6(address) {
			ipv4 = ""
			ipv6 = address
		} else {
			ipv4 = ""
			ipv6 = ""
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "Invalid remote address"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusBadRequest, payload)
			return
		}
		break
	default:
		ipv4 = ""
		ipv6 = ""
		log.Printf("Invalid ingest mode configuration for key item name \"%s\"", keyItem.Name)
		payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "Server configuration error"}}}
		log.Printf("%+v", payload)
		c.JSON(http.StatusInternalServerError, payload)
		return
	}

	// Clean up previously created DNS resource records
	if ipv4 != "" || ipv6 != "" {
		log.Print("Cleaning up any previously created IPv4 resource records")
		if err := C.PowerDNS.DeleteIPv4ResourceRecord(keyItem.HostName); err != nil {
			log.Printf("%+v", err)
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv4 record deletion failed"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusInternalServerError, payload)
			return
		}
		log.Print("Cleaning up any previously created IPv6 resource records")
		if err := C.PowerDNS.DeleteIPv6ResourceRecord(keyItem.HostName); err != nil {
			log.Printf("%+v", err)
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv6 record deletion failed"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusInternalServerError, payload)
			return
		}
	}

	// Create new DNS resource records
	if ipv4 != "" {
		log.Print("Creating IPv4 resource records")
		if err := C.PowerDNS.AddIPv4ResourceRecord(keyItem.HostName, ipv4, keyItem.TTL); err != nil {
			log.Printf("%+v", err)
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv4 record creation failed"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusInternalServerError, payload)
			return
		}
	}
	if ipv6 != "" {
		log.Print("Creating IPv6 resource records")
		if err := C.PowerDNS.AddIPv6ResourceRecord(keyItem.HostName, ipv6, keyItem.TTL); err != nil {
			log.Printf("%+v", err)
			payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "IPv6 record creation failed"}}}
			log.Printf("%+v", payload)
			c.JSON(http.StatusInternalServerError, payload)
			return
		}
	}

	// Build response payload
	if keyItem.HostName != "" && keyItem.IngestMode != "" && (ipv4 != "" || ipv6 != "") {
		payload := HostSyncPayload{HostSyncObjects: []*HostSyncObject{{
			HostName:   keyItem.HostName,
			IngestMode: keyItem.IngestMode,
			IPv4:       ipv4,
			IPv6:       ipv6,
		}}}
		log.Printf("Updated \"%s\" successfully", keyItem.Name)
		c.JSON(http.StatusOK, payload)
		return
	}
	payload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: "HostSync request process error"}}}
	log.Printf("%+v", payload)
	c.JSON(http.StatusInternalServerError, payload)
	return
}
