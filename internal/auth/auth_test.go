package auth

import (
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"testing"
)

func TestCheckHost(t *testing.T) {
	t.Run("TestCheckHostEmptyReturnNotEmpty", func(t *testing.T) {
		if host, _ := checkHost(""); host != "" {
			t.Error("Empty host returns invalid host string")
		}
	})
	t.Run("TestCheckHostEmptyErrorNil", func(t *testing.T) {
		if _, err := checkHost(""); err == nil {
			t.Error("Empty host does not return error")
		}
	})
	t.Run("TestCheckTestHostReturnInvalid", func(t *testing.T) {
		if host, _ := checkHost("foo"); host != "foo" {
			t.Error("Test host returns invalid host string")
		}
	})
	t.Run("TestCheckTestHostReturnNotNil", func(t *testing.T) {
		if _, err := checkHost("foo"); err != nil {
			t.Error("Test host returns error")
		}
	})
}

func TestCheckKey(t *testing.T) {
	t.Run("TestCheckKeyEmptyReturnNotEmpty", func(t *testing.T) {
		if key, _ := checkKey(""); key != "" {
			t.Error("Empty key returns invalid key string")
		}
	})
	t.Run("TestCheckKeyEmptyErrorNil", func(t *testing.T) {
		if _, err := checkKey(""); err == nil {
			t.Error("Empty key does not return error")
		}
	})
	t.Run("TestCheckTestKeyReturnInvalid", func(t *testing.T) {
		if key, _ := checkKey("foo"); key != "foo" {
			t.Error("Test key returns invalid key string")
		}
	})
	t.Run("TestCheckTestKeyReturnNotNil", func(t *testing.T) {
		if _, err := checkKey("foo"); err != nil {
			t.Error("Test key returns error")
		}
	})
}

func TestCheckAuthorization(t *testing.T) {
	keyTable := []yamlconfig.Key{
		{
			Name:       "homeRouter",
			Enable:     true,
			Key:        "secret",
			HostName:   "home-router",
			IngestMode: yamlconfig.IngestModeGetParameter,
			TTL:        1337,
		},
	}

	t.Run("TestCheckAuthorizationEmptyReturnNotEmpty", func(t *testing.T) {
		if keyItem, _ := checkAuthorization(keyTable, "foo", "bar"); keyItem.Name != "" {
			t.Error("Invalid key name returns invalid key item")
		}
	})
	t.Run("TestCheckAuthorizationEmptyErrorNil", func(t *testing.T) {
		if _, err := checkAuthorization(keyTable, "foo", "bar"); err == nil {
			t.Error("Invalid key name does not return error")
		}
	})
	t.Run("TestCheckTestAuthorizationReturnInvalid", func(t *testing.T) {
		if keyItem, _ := checkAuthorization(keyTable, "homeRouter", "secret"); keyItem.Name != "homeRouter" {
			t.Error("Test key name returns invalid key item")
		}
	})
	t.Run("TestCheckTestAuthorizationReturnNotNil", func(t *testing.T) {
		if _, err := checkAuthorization(keyTable, "homeRouter", "secret"); err != nil {
			t.Error("Test key name returns error")
		}
	})
}
