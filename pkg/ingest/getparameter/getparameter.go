package getparameter

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/ginResponse"
	"github.com/joeig/dyndns-pdns/pkg/ingest"
	"log"
	"net/http"
)

type GetParameter struct {
	Ctx *gin.Context
}

func (g *GetParameter) Process() (*ingest.IPSet, error) {
	ipSet := &ingest.IPSet{
		IPv4: g.Ctx.Query("ipv4"),
		IPv6: g.Ctx.Query("ipv6"),
	}

	log.Printf("Received ipv4=\"%s\" ipv6=\"%s\"", ipSet.IPv4, ipSet.IPv6)

	if !ipSet.HasIPv4() && !ipSet.HasIPv6() {
		ginResponse.GinJSONError(g.Ctx, http.StatusBadRequest, "IPv4 as well as IPv6 parameter missing")
		return ipSet, &ingest.Error{}
	}

	if ipSet.HasIPv4() && !ipSet.IsIPv4() {
		ginResponse.GinJSONError(g.Ctx, http.StatusBadRequest, "IPv4 address invalid")
		return ipSet, &ingest.Error{}
	}

	if ipSet.HasIPv6() && !ipSet.IsIPv6() {
		ginResponse.GinJSONError(g.Ctx, http.StatusBadRequest, "IPv6 address invalid")
		return ipSet, &ingest.Error{}
	}

	return ipSet, nil
}
