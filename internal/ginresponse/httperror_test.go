package ginresponse

import (
	"testing"
)

func TestError(t *testing.T) {
	e := &HTTPError{Message: "Foo"}
	if e.Error() != "Foo" {
		t.Error("HTTPError method call returns invalid value")
	}
}
