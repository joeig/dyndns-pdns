package getparameter

import (
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
)

type GetParameter struct {
	IPv4 string
	IPv6 string
}

func (g *GetParameter) Process() (*ingest.IPSet, error) {
	ipSet := &ingest.IPSet{
		IPv4: g.IPv4,
		IPv6: g.IPv6,
	}

	if !ipSet.HasIPv4() && !ipSet.HasIPv6() {
		return ipSet, &genericerror.GenericError{Message: "IPv4 as well as IPv6 parameter missing"}
	}

	if ipSet.HasIPv4() && !ipSet.IsIPv4() {
		return ipSet, &genericerror.GenericError{Message: "IPv4 address invalid"}
	}

	if ipSet.HasIPv6() && !ipSet.IsIPv6() {
		return ipSet, &genericerror.GenericError{Message: "IPv6 address invalid"}
	}

	return ipSet, nil
}
