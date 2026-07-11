package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
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
	configFile := "../../configs/config.test.yml"
	yamlconfig.C = yamlconfig.Config{}
	yamlconfig.ParseConfig(&yamlconfig.C, &configFile)
	yamlconfig.SetDNSProvider(&yamlconfig.ActiveDNSProvider)
	Dry = true
	yamlconfig.C.PowerDNS.Dry = Dry
	router := setupGinEngine()

	// OK
	t.Run("TestGetParameterIPv4IPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestGetParameterIPv4OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestGetParameterIPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv4OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/officeRouter/sync?key=topSecret", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv4PortOK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/officeRouter/sync?key=topSecret", "127.0.0.1:1337", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/officeRouter/sync?key=topSecret", "::1", http.StatusOK)
	})
	t.Run("TestRemoteAddressIPv6PortOK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/officeRouter/sync?key=topSecret", "[::1]:1337", http.StatusOK)
	})

	// Forbidden
	t.Run("TestUnknownDeviceNameForbidden", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/unknownDevice/sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusForbidden)
	})
	t.Run("TestInvalidKeyForbidden", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=wrongKey&ipv6=::1", "127.0.0.1", http.StatusForbidden)
	})

	// Unauthorized
	t.Run("TestMissingDeviceNameUnauthorized", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host//sync?key=secret&ipv6=::1", "127.0.0.1", http.StatusUnauthorized)
	})
	t.Run("TestMissingKeyUnauthorized", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync", "127.0.0.1", http.StatusUnauthorized)
	})

	// BadRequest
	t.Run("TestGetParameterMissingBadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestInvalidIPv4BadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=foo", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestInvalidIPv6BadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv6=foo", "127.0.0.1", http.StatusBadRequest)
	})

	// Response headers
	t.Run("TestCacheControl", func(t *testing.T) {
		res := assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
		if res.Header().Get("Cache-Control") == "" {
			t.Errorf("Cache-Control is missing")
		}
	})
	t.Run("TestRequestID", func(t *testing.T) {
		res := assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&ipv6=::1", "127.0.0.1", http.StatusOK)
		if res.Header().Get("X-Request-ID") == "" {
			t.Errorf("X-Request-ID is missing")
		}
	})
}

func TestHostSyncWithPrefix(t *testing.T) {
	configFile := "../../configs/config.test.with-dynamic.yml"
	yamlconfig.C = yamlconfig.Config{}
	yamlconfig.ParseConfig(&yamlconfig.C, &configFile)
	yamlconfig.SetDNSProvider(&yamlconfig.ActiveDNSProvider)
	Dry = true
	yamlconfig.C.PowerDNS.Dry = Dry
	router := setupGinEngine()

	// Prefix functionality tests with dynamic keys
	t.Run("TestPrefixOnlyOK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&prefix=2001:db8:1234::/48", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestPrefixWithIPv4OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv4=127.0.0.1&prefix=2001:db8:1234::/48", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestPrefixWithIPv6OK", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&ipv6=::1&prefix=2001:db8:1234::/48", "127.0.0.1", http.StatusOK)
	})
	t.Run("TestInvalidPrefixBadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret&prefix=invalid-prefix", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestMissingPrefixBadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/homeRouter/sync?key=secret", "127.0.0.1", http.StatusBadRequest)
	})
	t.Run("TestMissingDynamicKeyBadRequest", func(t *testing.T) {
		assertHostSyncComponent(t, router, http.MethodGet, "/v1/host/dumbRouter/sync?key=secret&prefix=2001:db8:1234::/48", "127.0.0.1", http.StatusBadRequest)
	})
}

func TestCombinePrefixWithInterfaceID(t *testing.T) {
	// Valid combinations
	t.Run("ValidCIDRPrefixWithInterfaceID", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("2001:db8:1234::/48", "::1")
		if err != nil {
			t.Errorf("should not fail, but returned error: %v", err)
		}
		if got != "2001:db8:1234::1" {
			t.Errorf("returned %v instead of %v", got, "2001:db8:1234::1")
		}
	})
	t.Run("ValidIPPrefixWithInterfaceID", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("2001:db8:1234::", "::1")
		if err != nil {
			t.Errorf("should not fail, but returned error: %v", err)
		}
		if got != "2001:db8:1234::1" {
			t.Errorf("returned %v instead of %v", got, "2001:db8:1234::1")
		}
	})
	t.Run("DifferentInterfaceID", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("2001:db8:1234::/48", "::abcd:ef01")
		if err != nil {
			t.Errorf("should not fail, but returned error: %v", err)
		}
		if got != "2001:db8:1234::abcd:ef01" {
			t.Errorf("returned %v instead of %v", got, "2001:db8:1234::abcd:ef01")
		}
	})

	// Error cases
	t.Run("InvalidPrefix", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("invalid-prefix", "::1")
		if err == nil {
			t.Errorf("should fail, but returned no error")
		}
		if got != "" {
			t.Errorf("returned %v instead of empty string", got)
		}
	})
	t.Run("InvalidInterfaceID", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("2001:db8:1234::/48", "invalid-interface-id")
		if err == nil {
			t.Errorf("should fail, but returned no error")
		}
		if got != "" {
			t.Errorf("returned %v instead of empty string", got)
		}
	})
	t.Run("IPv4Prefix", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("192.168.1.0/24", "::1")
		if err == nil {
			t.Errorf("should fail, but returned no error")
		}
		if got != "" {
			t.Errorf("returned %v instead of empty string", got)
		}
	})
	t.Run("IPv4InterfaceID", func(t *testing.T) {
		got, err := combinePrefixWithInterfaceID("2001:db8:1234::/48", "192.168.1.1")
		if err == nil {
			t.Errorf("should fail, but returned no error")
		}
		if got != "" {
			t.Errorf("returned %v instead of empty string", got)
		}
	})
}
