package remoteaddress

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/ginResponse"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"log"
	"net"
	"net/http"
)

type RemoteAddress struct {
	Ctx *gin.Context
}

func isolateHostAddress(remoteAddress string) string {
	// Under certain circumstances, RemoteAddr contains also the port number
	address, _, err := net.SplitHostPort(remoteAddress)
	if err != nil {
		return remoteAddress
	}
	return address
}

func (r *RemoteAddress) Process() (*ingest.IPSet, error) {
	address := isolateHostAddress(r.Ctx.Request.RemoteAddr)
	log.Printf("Received address=\"%s\"", address)

	ipSet := &ingest.IPSet{}

	if govalidator.IsIPv4(address) {
		ipSet.IPv4 = address
	} else if govalidator.IsIPv6(address) {
		ipSet.IPv6 = address
	} else {
		ginResponse.GinJSONError(r.Ctx, http.StatusBadRequest, "Invalid remote address")
		return ipSet, &ingest.Error{}
	}

	return ipSet, nil
}
