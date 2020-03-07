package genericerror

import (
	"testing"
)

func TestError(t *testing.T) {
	e := &GenericError{"Foo"}
	if e.Error() != "Foo" {
		t.Error("GenericError method call returns invalid value")
	}
}
