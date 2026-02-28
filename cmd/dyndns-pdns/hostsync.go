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
	"net"
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
		if err := cleanUpOutdatedResourceRecords(ipSet, keyItem, keyItem.HostName); err != nil {
			ginresponse.GinJSONError(ctx, err)
			return
		}
	}

	hostSyncObjects := []*HostSyncObject{}
	hostSyncObject, err := createNewResourceRecords(ipSet, keyItem)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}
	if hostSyncObject != nil {
		hostSyncObjects = append(hostSyncObjects, hostSyncObject)
	}

	prefix, err := getPrefix(ctx, keyItem)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}
	prefixhostSyncObjects, err := updatePrefixes(keyItem, prefix)
	if err != nil {
		ginresponse.GinJSONError(ctx, err)
		return
	}
	hostSyncObjects = append(hostSyncObjects, prefixhostSyncObjects...)

	buildResponsePayload(ctx, keyItem, hostSyncObjects)
}

// Combine a prefix (CIDR with prefix length or IP address with /64 fallback) with an interface ID
func combinePrefixWithInterfaceID(prefix string, interfaceID string) (string, error) {
	_, prefixNet, err := net.ParseCIDR(prefix)
	if err != nil {
		prefixIP := net.ParseIP(prefix)
		if prefixIP == nil || len(prefixIP) != net.IPv6len {
			return "", &ginresponse.HTTPError{Message: "Invalid IPv6 prefix: " + prefix, HTTPErrorCode: http.StatusBadRequest}
		}
		prefixNet = &net.IPNet{
			IP:   prefixIP,
			Mask: net.CIDRMask(64, 128),
		}
	}
	if prefixNet.IP.To4() != nil {
		return "", &ginresponse.HTTPError{Message: "IPv4 prefixes are not supported: " + prefix, HTTPErrorCode: http.StatusBadRequest}
	}

	interfaceIDIP := net.ParseIP(interfaceID)
	if interfaceIDIP == nil || interfaceIDIP.To4() != nil {
		return "", &ginresponse.HTTPError{Message: "Invalid IPv6 interface ID: " + interfaceID, HTTPErrorCode: http.StatusBadRequest}
	}

	combined := make(net.IP, net.IPv6len)
	for i := range net.IPv6len {
		combined[i] = (prefixNet.IP[i] & prefixNet.Mask[i]) | (interfaceIDIP[i] & ^prefixNet.Mask[i])
	}
	return combined.String(), nil
}

func updatePrefixes(keyItem *yamlconfig.Key, prefix string) ([]*HostSyncObject, error) {
	var hostSyncObjects []*HostSyncObject

	if len(keyItem.DynamicKeys) == 0 {
		return hostSyncObjects, nil
	}

	if prefix == "" {
		return nil, &ginresponse.HTTPError{Message: "No prefix provided (updating prefixes only works in getParameter mode)", HTTPErrorCode: http.StatusBadRequest}
	}

	for _, dynamicHost := range keyItem.DynamicKeys {
		if !dynamicHost.Enable {
			continue
		}

		if dynamicHost.InterfaceID == "" {
			return nil, &ginresponse.HTTPError{Message: "No interfaceID configured for host: " + dynamicHost.Name, HTTPErrorCode: http.StatusBadRequest}
		}

		combinedIPv6, err := combinePrefixWithInterfaceID(prefix, dynamicHost.InterfaceID)
		if err != nil {
			return nil, err
		}

		prefixIPSet := &ingest.IPSet{IPv6: combinedIPv6}

		if err := cleanUpOutdatedResourceRecords(prefixIPSet, keyItem, dynamicHost.HostName); err != nil {
			return nil, err
		}

		if err := createNewIPv6ResourceRecord(prefixIPSet, keyItem, dynamicHost.HostName); err != nil {
			return nil, err
		}

		hostSyncObjects = append(hostSyncObjects, &HostSyncObject{
			HostName:    dynamicHost.HostName,
			IngestMode:  keyItem.IngestMode,
			CleanUpMode: keyItem.CleanUpMode,
			TTL:         int(keyItem.TTL),
			IPv6:        combinedIPv6,
		})
	}

	return hostSyncObjects, nil
}

func getIngestModeHandler(ctx *gin.Context, desiredIngestModeType ingest.ModeType) (ingest.Mode, error) {
	var activeIngestMode ingest.Mode

	switch desiredIngestModeType {
	case yamlconfig.IngestModeGetParameter:
		ipv4 := ctx.Query("ipv4")
		ipv6 := ctx.Query("ipv6")
		prefix := ctx.Query("prefix")
		log.Printf("Received ipv4=\"%s\" ipv6=\"%s\" prefix=\"%s\"", ipv4, ipv6, prefix)

		activeIngestMode = &getparameter.GetParameter{IPv4: ipv4, IPv6: ipv6, Prefix: prefix}

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
	ipSet, err := activeIngestMode.GetIPSet()
	if err != nil {
		return ipSet, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusBadRequest}
	}

	log.Printf("Gathered ipSet: %+v", ipSet)
	return ipSet, nil
}

func getPrefix(ctx *gin.Context, keyItem *yamlconfig.Key) (string, error) {
	activeIngestMode, err := getIngestModeHandler(ctx, keyItem.IngestMode)
	if err != nil {
		log.Printf("Unable to initialise ingests mode for \"%s\": %s", keyItem.Name, err.Error())
		return "", err
	}
	log.Printf("Processing ingest for %+v mode", keyItem.IngestMode)
	return activeIngestMode.GetPrefix()
}

func buildResponsePayload(ctx *gin.Context, keyItem *yamlconfig.Key, hostSyncObjects []*HostSyncObject) {
	payload := HostSyncPayload{HostSyncObjects: hostSyncObjects}
	log.Printf("Updated \"%s\" successfully", keyItem.Name)
	ctx.JSON(http.StatusOK, payload)
}
