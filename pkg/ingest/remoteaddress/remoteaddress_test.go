package remoteaddress

import "testing"

func TestIsolateHostAddress(t *testing.T) {
	if isolateHostAddress("foo") != "foo" {
		t.Error("Invalid host address returned")
	}
	if isolateHostAddress("foo:80") != "foo" {
		t.Error("Port was not stripped successfully")
	}
}
