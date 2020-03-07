package remoteaddress

import (
	"github.com/asaskevich/govalidator"
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"net"
)

type RemoteAddress struct {
	Address string
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
	address := isolateHostAddress(r.Address)

	ipSet := &ingest.IPSet{}

	if govalidator.IsIPv4(address) {
		ipSet.IPv4 = address
	} else if govalidator.IsIPv6(address) {
		ipSet.IPv6 = address
	} else {
		return ipSet, &genericerror.GenericError{Message: "Invalid remote address"}
	}

	return ipSet, nil
}
