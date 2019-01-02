package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	router := getGinEngine()
	res := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/v1/health", nil)
	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("HTTP request does not return %v", http.StatusOK)
	}
}
