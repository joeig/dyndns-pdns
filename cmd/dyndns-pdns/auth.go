package main

import "github.com/joeig/dyndns-pdns/internal/genericerror"

func checkHost(host string) (string, error) {
	if host == "" {
		return "", &genericerror.GenericError{Message: "Host parameter missing"}
	}
	return host, nil
}

func checkKey(key string) (string, error) {
	if key == "" {
		return "", &genericerror.GenericError{Message: "Key parameter missing"}
	}
	return key, nil
}

func checkAuthorization(keyTable []Key, name string, key string) (*Key, error) {
	for _, keyItem := range keyTable {
		if keyItem.Enable && keyItem.Name == name && keyItem.Key == key {
			return &keyItem, nil
		}
	}
	return &Key{}, &genericerror.GenericError{Message: "Permission denied"}
}
