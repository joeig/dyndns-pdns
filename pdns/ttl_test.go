package pdns

import "testing"

func TestGetTTL(t *testing.T) {
	t.Run("TestGetTTLNotDefault", func(t *testing.T) {
		if getTTL(0, 10) != 10 {
			t.Error("Zero key item TTL returns invalid default TTL")
		}
	})
	t.Run("TestGetTTLInvalid", func(t *testing.T) {
		if getTTL(1337, 10) != 1337 {
			t.Error("Test key item TTL returns invalid TTL")
		}
	})
}
