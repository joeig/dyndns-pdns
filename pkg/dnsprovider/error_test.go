package dnsprovider

import (
	"testing"
)

func TestError(t *testing.T) {
	e := &Error{"Foo"}
	if e.Error() != "Foo" {
		t.Error("Error method call returns invalid value")
	}
}
