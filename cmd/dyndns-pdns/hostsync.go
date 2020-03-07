package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"github.com/joeig/dyndns-pdns/internal/ginresponse"
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
	HostName    string          `json:"hostName"`
	IngestMode  ingest.ModeType `json:"ingestMode"`
	CleanUpMode CleanUpModeType `json:"cleanUpMode"`
	TTL         int             `json:"ttl"`
	IPv4        string          `json:"ipv4"`
	IPv6        string          `json:"ipv6"`
}

// HostSync Gin route
func HostSync(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache")

	name, err := getName(ctx)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	key, err := getKey(ctx)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	keyItem, err := getKeyItem(name, key)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	ipSet, err := getIPAddresses(ctx, keyItem)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	if ipSet.HasIPv4() || ipSet.HasIPv6() {
		if err := cleanUpOutdatedResourceRecords(ipSet, keyItem); err != nil {
			ginresponse.GinJSONError(ctx, err)
			return
		}
	}

	if err := createNewResourceRecords(ipSet, keyItem); err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	_ = buildResponsePayload(ctx, keyItem, ipSet)
}

func getName(ctx *gin.Context) (string, error) {
	name, err := checkHost(ctx.Param("name"))
	if err != nil {
		return name, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusUnauthorized}
	}

	log.Printf("Received name=\"%s\"", name)
	return name, nil
}

func getKey(ctx *gin.Context) (string, error) {
	key, err := checkKey(ctx.Query("key"))
	if err != nil {
		return key, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusUnauthorized}
	}

	log.Printf("Received key=\"%s\"", key)
	return key, nil
}

func getKeyItem(name string, key string) (*Key, error) {
	keyItem, err := checkAuthorization(C.KeyTable, name, key)
	if err != nil {
		return keyItem, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusForbidden}
	}

	log.Printf("Found key item: %+v", keyItem)
	return keyItem, nil
}

func getIngestModeHandler(ctx *gin.Context, desiredIngestModeType ingest.ModeType) (ingest.Mode, error) {
	var activeIngestMode ingest.Mode

	switch desiredIngestModeType {
	case IngestModeGetParameter:
		ipv4 := ctx.Query("ipv4")
		ipv6 := ctx.Query("ipv6")
		log.Printf("Received ipv4=\"%s\" ipv6=\"%s\"", ipv4, ipv6)

		activeIngestMode = &getparameter.GetParameter{IPv4: ipv4, IPv6: ipv6}

	case IngestModeRemoteAddress:
		address := ctx.Request.RemoteAddr
		log.Printf("Received address=\"%s\"", address)

		activeIngestMode = &remoteaddress.RemoteAddress{Address: ctx.Request.RemoteAddr}

	default:
		return activeIngestMode, &ginresponse.HTTPError{Message: "Server configuration error: Invalid ingest mode", HTTPErrorCode: http.StatusBadRequest}
	}

	return activeIngestMode, nil
}

func getIPAddresses(ctx *gin.Context, keyItem *Key) (*ingest.IPSet, error) {
	activeIngestMode, err := getIngestModeHandler(ctx, keyItem.IngestMode)
	if err != nil {
		log.Printf("Unable to initialise ingests mode for \"%s\": %s", keyItem.Name, err.Error())
		return &ingest.IPSet{}, err
	}

	log.Printf("Processing ingest for %+v mode", keyItem.IngestMode)
	ipSet, err := activeIngestMode.Process()
	if err != nil {
		return ipSet, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusBadRequest}
	}

	log.Printf("Gathered ipSet: %+v", ipSet)
	return ipSet, nil
}

func cleanUpOutdatedResourceRecords(ipSet *ingest.IPSet, keyItem *Key) error {
	if keyItem.CleanUpMode == CleanUpModeAny || (keyItem.CleanUpMode == CleanUpModeRequestBased && ipSet.HasIPv4()) {
		log.Print("Cleaning up any previously created IPv4 resource records")

		if err := activeDNSProvider.DeleteIPv4ResourceRecord(keyItem.HostName); err != nil {
			log.Printf("%+v", err)
			return &ginresponse.HTTPError{Message: "IPv4 record deletion failed", HTTPErrorCode: http.StatusInternalServerError}
		}
	} else {
		log.Print("Skipping clean up of previously created IPv4 resource records")
	}

	if keyItem.CleanUpMode == CleanUpModeAny || (keyItem.CleanUpMode == CleanUpModeRequestBased && ipSet.HasIPv6()) {
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

func createNewIPv4ResourceRecord(ipSet *ingest.IPSet, keyItem *Key) error {
	log.Print("Creating IPv4 resource records")

	if err := activeDNSProvider.AddIPv4ResourceRecord(keyItem.HostName, ipSet.IPv4, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv4 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewIPv6ResourceRecord(ipSet *ingest.IPSet, keyItem *Key) error {
	log.Print("Creating IPv6 resource records")

	if err := activeDNSProvider.AddIPv6ResourceRecord(keyItem.HostName, ipSet.IPv6, keyItem.TTL); err != nil {
		log.Printf("%+v", err)
		return &ginresponse.HTTPError{Message: "IPv6 record creation failed", HTTPErrorCode: http.StatusInternalServerError}
	}

	return nil
}

func createNewResourceRecords(ipSet *ingest.IPSet, keyItem *Key) error {
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

func buildResponsePayload(ctx *gin.Context, keyItem *Key, ipSet *ingest.IPSet) error {
	if keyItem.HostName != "" && keyItem.IngestMode != "" && (ipSet.HasIPv4() || ipSet.HasIPv6()) {
		payload := HostSyncPayload{HostSyncObjects: []*HostSyncObject{{
			HostName:    keyItem.HostName,
			IngestMode:  keyItem.IngestMode,
			CleanUpMode: keyItem.CleanUpMode,
			IPv4:        ipSet.IPv4,
			IPv6:        ipSet.IPv6,
		}}}
		log.Printf("Updated \"%s\" successfully", keyItem.Name)
		ctx.JSON(http.StatusOK, payload)
		return nil
	}

	ginresponse.GinJSONError(ctx, &ginresponse.HTTPError{Message: "HostSync request processing error", HTTPErrorCode: http.StatusInternalServerError})
	return &genericerror.GenericError{}
}
