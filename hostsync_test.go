package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertHostSyncComponent(t *testing.T, router *gin.Engine, method string, url string, remoteAddr string, assertedCode int) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	req.RemoteAddr = remoteAddr
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != assertedCode {
		t.Errorf("HTTP request to \"%s\" returned %d instead of %d", url, res.Code, assertedCode)
	}
	return res
}

func TestHostSync(t *testing.T) {
	configFile := "config.dist.yml"
	parseConfig(&C, &configFile)
	Dry = true
	C.PowerDNS.Dry = Dry
	router := getGinEngine()

	// OK
	t.Run("TestGetParameterIPv4IPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestGetParameterIPv4OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestGetParameterIPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv4OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/officeRouter/sync?key=topSecret", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv4PortOK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/officeRouter/sync?key=topSecret", "127.0.0.1:1337", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/officeRouter/sync?key=topSecret", "::1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv6PortOK", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/officeRouter/sync?key=topSecret", "[::1]:1337", http.StatusOK)
	})

	// Forbidden
	t.Run("TestUnknownDeviceNameForbidden", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/unknownDevice/sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusForbidden)
	})
	t.Run("TestInvalidKeyForbidden", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=wrongKey&ipv6=::1", "127.0.0.1", http.StatusForbidden)
	})

	// Unauthorized
	t.Run("TestMissingDeviceNameUnauthorized", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host//sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusUnauthorized)
	})
	t.Run("TestMissingKeyUnauthorized", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync", "127.0.0.1", http.StatusUnauthorized)
	})

	// BadRequest
	t.Run("TestGetParameterMissingBadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestInvalidIPv4BadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv4=foo", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestInvalidIPv6BadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv6=foo", "127.0.0.1", http.StatusBadRequest)
	})

	// Response headers
	t.Run("TestCacheControl", func(t *testing.T) {
		res := assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
		if res.HeaderMap.Get("Cache-Control") == "" {
			t.Errorf("Cache-Control is missing")
		}
	})
	t.Run("TestRequestID", func(t *testing.T) {
		res := assertHostSyncComponent(t, router, "GET", "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
		if res.HeaderMap.Get("X-Request-ID") == "" {
			t.Errorf("X-Request-ID is missing")
		}
	})
}
