package main

import (
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"github.com/joeig/dyndns-pdns/internal/ginresponse"
	"github.com/joeig/dyndns-pdns/internal/yamlconfig"
	"log"
	"net/http"
)

func checkHost(host string) (string, error) {
	if host == "" {
		return "", &genericerror.GenericError{Message: "Host parameter missing"}
	}
	return host, nil
}

func getName(host string) (string, error) {
	name, err := checkHost(host)
	if err != nil {
		return name, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusUnauthorized}
	}

	log.Printf("Received name=\"%s\"", name)
	return name, nil
}

func checkKey(key string) (string, error) {
	if key == "" {
		return "", &genericerror.GenericError{Message: "Key parameter missing"}
	}
	return key, nil
}

func getKey(key string) (string, error) {
	key, err := checkKey(key)
	if err != nil {
		return key, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusUnauthorized}
	}

	log.Printf("Received key=\"%s\"", key)
	return key, nil
}

func checkAuthorization(keyTable []yamlconfig.Key, name string, key string) (*yamlconfig.Key, error) {
	for _, keyItem := range keyTable {
		if keyItem.Enable && keyItem.Name == name && keyItem.Key == key {
			return &keyItem, nil
		}
	}
	return &yamlconfig.Key{}, &genericerror.GenericError{Message: "Permission denied"}
}

func getKeyItem(name string, key string) (*yamlconfig.Key, error) {
	keyItem, err := checkAuthorization(yamlconfig.C.KeyTable, name, key)
	if err != nil {
		return keyItem, &ginresponse.HTTPError{Message: err.Error(), HTTPErrorCode: http.StatusForbidden}
	}

	log.Printf("Found key item: %+v", keyItem)
	return keyItem, nil
}
