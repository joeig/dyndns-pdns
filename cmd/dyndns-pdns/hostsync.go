package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/ginResponse"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"github.com/joeig/dyndns-pdns/pkg/ingest/getparameter"
	"github.com/joeig/dyndns-pdns/pkg/ingest/remoteaddress"
	"log"
	"net/http"
)

// HostSyncPayload returns a payload containing HostSyncObjects
type HostSyncPayload struct {
	HostSyncObjects []*HostSyncObject `json:"hostSyncObjects"`
}

// HostSyncObject contains a payload for the requester in order to identify the values that have been stored
type HostSyncObject struct {
	HostName   string                `json:"hostName"`
	IngestMode ingest.IngestModeType `json:"ingestMode"`
	TTL        int                   `json:"ttl"`
	IPv4       string                `json:"ipv4"`
	IPv6       string                `json:"ipv6"`
}

// HostSync Gin route
func HostSync(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache")

	name, err := getName(ctx)
	if err != nil {
		return
	}

	key, err := getKey(ctx)
	if err != nil {
		return
	}

	keyItem, err := getKeyItem(ctx, name, key)
	if err != nil {
		return
	}

	ipSet, err := getIPAddresses(ctx, keyItem)
	if err != nil {
		return
	}

	if ipSet.HasIPv4() || ipSet.HasIPv6() {
		if cleanUpOutdatedResourceRecords(ctx, keyItem.HostName) != nil {
			return
		}
	}

	if createNewResourceRecords(ctx, ipSet, keyItem) != nil {
		return
	}

	_ = buildResponsePayload(ctx, keyItem, ipSet)
}

func getName(ctx *gin.Context) (string, error) {
	name, err := checkHost(ctx.Param("name"))
	log.Printf("Received name=\"%s\"", name)
	if err != nil {
		ginResponse.GinJSONError(ctx, http.StatusUnauthorized, err.Error())
		return name, err
	}

	return name, nil
}

func getKey(ctx *gin.Context) (string, error) {
	key, err := checkKey(ctx.Query("key"))
	log.Printf("Received key=\"%s\"", key)
	if err != nil {
		ginResponse.GinJSONError(ctx, http.StatusUnauthorized, err.Error())
		return key, err
	}

	return key, nil
}

func getKeyItem(ctx *gin.Context, name string, key string) (*Key, error) {
	keyItem, err := checkAuthorization(C.KeyTable, name, key)
	if err != nil {
		ginResponse.GinJSONError(ctx, http.StatusForbidden, err.Error())
		return keyItem, err
	}

	log.Printf("Found key item: %+v", keyItem)
	return keyItem, nil
}

func getIngestModeHandler(ctx *gin.Context, desiredIngestModeType ingest.IngestModeType) (ingest.IngestMode, error) {
	var activeIngestMode ingest.IngestMode

	switch desiredIngestModeType {
	case IngestModeGetParameter:
		activeIngestMode = &getparameter.GetParameter{Ctx: ctx}
	case IngestModeRemoteAddress:
		activeIngestMode = &remoteaddress.RemoteAddress{Ctx: ctx}
	default:
		ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "Server configuration error")
		return activeIngestMode, &Error{}
	}

	return activeIngestMode, nil
}

func getIPAddresses(ctx *gin.Context, keyItem *Key) (*ingest.IPSet, error) {
	activeIngestMode, err := getIngestModeHandler(ctx, keyItem.IngestMode)
	if err != nil {
		log.Printf("Invalid ingest mode configuration for key item name \"%s\"", keyItem.Name)
		return &ingest.IPSet{}, err
	}

	log.Printf("Processing ingest for %+v mode", keyItem.IngestMode)
	return activeIngestMode.Process()
}

func cleanUpOutdatedResourceRecords(ctx *gin.Context, hostname string) error {
	log.Print("Cleaning up any previously created IPv4 resource records")

	if err := activeDNSProvider.DeleteIPv4ResourceRecord(hostname); err != nil {
		log.Printf("%+v", err)
		ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "IPv4 record deletion failed")
		return &Error{}
	}

	log.Print("Cleaning up any previously created IPv6 resource records")

	if err := activeDNSProvider.DeleteIPv6ResourceRecord(hostname); err != nil {
		log.Printf("%+v", err)
		ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "IPv6 record deletion failed")
		return &Error{}
	}

	return nil
}

func createNewIPv4ResourceRecord(ctx *gin.Context, ipSet *ingest.IPSet, keyItem *Key) error {
	log.Print("Creating IPv4 resource records")

	if err := activeDNSProvider.AddIPv4ResourceRecord(keyItem.HostName, ipSet.IPv4, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "IPv4 record creation failed")
		return &Error{}
	}

	return nil
}

func createNewIPv6ResourceRecord(ctx *gin.Context, ipSet *ingest.IPSet, keyItem *Key) error {
	log.Print("Creating IPv6 resource records")

	if err := activeDNSProvider.AddIPv6ResourceRecord(keyItem.HostName, ipSet.IPv6, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "IPv6 record creation failed")
		return &Error{}
	}

	return nil
}

func createNewResourceRecords(ctx *gin.Context, ipSet *ingest.IPSet, keyItem *Key) error {
	if ipSet.HasIPv4() {
		if err := createNewIPv4ResourceRecord(ctx, ipSet, keyItem); err != nil {
			return err
		}
	}

	if ipSet.HasIPv6() {
		if err := createNewIPv6ResourceRecord(ctx, ipSet, keyItem); err != nil {
			return err
		}
	}

	return nil
}

func buildResponsePayload(ctx *gin.Context, keyItem *Key, ipSet *ingest.IPSet) error {
	if keyItem.HostName != "" && keyItem.IngestMode != "" && (ipSet.HasIPv4() || ipSet.HasIPv6()) {
		payload := HostSyncPayload{HostSyncObjects: []*HostSyncObject{{
			HostName:   keyItem.HostName,
			IngestMode: keyItem.IngestMode,
			IPv4:       ipSet.IPv4,
			IPv6:       ipSet.IPv6,
		}}}
		log.Printf("Updated \"%s\" successfully", keyItem.Name)
		ctx.JSON(http.StatusOK, payload)
		return nil
	}

	ginResponse.GinJSONError(ctx, http.StatusInternalServerError, "HostSync request processing error")
	return &Error{}
}
