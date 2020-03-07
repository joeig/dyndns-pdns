package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/auth"
	"github.com/joeig/dyndns-pdns/internal/ginresponse"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
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
	HostName    string                     `json:"hostName"`
	IngestMode  ingest.ModeType            `json:"ingestMode"`
	CleanUpMode yamlconfig.CleanUpModeType `json:"cleanUpMode"`
	TTL         int                        `json:"ttl"`
	IPv4        string                     `json:"ipv4"`
	IPv6        string                     `json:"ipv6"`
}

// HostSync Gin route
func HostSync(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache")

	name, err := auth.GetName(ctx.Param("name"))
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	key, err := auth.GetKey(ctx.Query("key"))
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}

	keyItem, err := auth.GetKeyItem(name, key)
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

	buildResponsePayload(ctx, keyItem, ipSet)
}

func getIngestModeHandler(ctx *gin.Context, desiredIngestModeType ingest.ModeType) (ingest.Mode, error) {
	var activeIngestMode ingest.Mode

	switch desiredIngestModeType {
	case yamlconfig.IngestModeGetParameter:
		ipv4 := ctx.Query("ipv4")
		ipv6 := ctx.Query("ipv6")
		log.Printf("Received ipv4=\"%s\" ipv6=\"%s\"", ipv4, ipv6)

		activeIngestMode = &getparameter.GetParameter{IPv4: ipv4, IPv6: ipv6}

	case yamlconfig.IngestModeRemoteAddress:
		address := ctx.Request.RemoteAddr
		log.Printf("Received address=\"%s\"", address)

		activeIngestMode = &remoteaddress.RemoteAddress{Address: ctx.Request.RemoteAddr}

	default:
		return activeIngestMode, &ginresponse.HTTPError{Message: "Server configuration error: Invalid ingest mode", HTTPErrorCode: http.StatusBadRequest}
	}

	return activeIngestMode, nil
}

func getIPAddresses(ctx *gin.Context, keyItem *yamlconfig.Key) (*ingest.IPSet, error) {
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

func buildResponsePayload(ctx *gin.Context, keyItem *yamlconfig.Key, ipSet *ingest.IPSet) {
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
		return
	}

	ginresponse.GinJSONError(ctx, &ginresponse.HTTPError{Message: "HostSync request processing error", HTTPErrorCode: http.StatusInternalServerError})
}
