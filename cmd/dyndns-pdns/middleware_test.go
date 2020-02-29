package main

import "testing"

func TestGenerateRequestID(t *testing.T) {
	requestID, err := generateRequestID()
	if requestID == "" && err == nil {
		t.Error("Request ID was not generated properly, but error is nil")
	}
}
