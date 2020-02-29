package ingest

import "testing"

func TestIPSetHasIPv4(t *testing.T) {
	ipSet1 := &IPSet{
		IPv4: "foo",
		IPv6: "bar",
	}
	if !ipSet1.HasIPv4() {
		t.Error("IPSet has IPv4, but returns false")
	}
	ipSet2 := &IPSet{
		IPv6: "bar",
	}
	if ipSet2.HasIPv4() {
		t.Error("IPSet has no IPv4, but returns true")
	}
}

func TestIPSetHasIPv6(t *testing.T) {
	ipSet1 := &IPSet{
		IPv4: "foo",
		IPv6: "bar",
	}
	if !ipSet1.HasIPv6() {
		t.Error("IPSet has IPv6, but returns false")
	}
	ipSet2 := &IPSet{
		IPv4: "bar",
	}
	if ipSet2.HasIPv6() {
		t.Error("IPSet has no IPv6, but returns true")
	}
}
