package ginresponse

import (
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"net/http"
	"testing"
)

func generateTestHTTPError(message string, httpErrorCode int) error {
	return &HTTPError{Message: message, HTTPErrorCode: httpErrorCode}
}

func generateTestGenericError(message string) error {
	return &genericerror.GenericError{Message: message}
}

func TestGatherHTTPErrorCode(t *testing.T) {
	testError := generateTestHTTPError("", 0)
	if gatherHTTPErrorCode(&testError) != DefaultHTTPErrorCode {
		t.Error("Error code is not 500, even though input error code is zero")
	}

	testError = generateTestHTTPError("", http.StatusBadRequest)
	if gatherHTTPErrorCode(&testError) != http.StatusBadRequest {
		t.Error("Specified error code is not returned")
	}

	testError = generateTestGenericError("")
	if gatherHTTPErrorCode(&testError) != DefaultHTTPErrorCode {
		t.Error("Generic error code is not returned")
	}
}
