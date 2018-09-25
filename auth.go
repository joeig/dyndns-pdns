package main

func checkHost(host string) (string, error) {
	if host == "" {
		return "", &Error{"Host parameter missing"}
	}
	return host, nil
}

func checkKey(key string) (string, error) {
	if key == "" {
		return "", &Error{"Key parameter missing"}
	}
	return key, nil
}

func checkAuthorization(keyTable []Key, name string, key string) (Key, error) {
	for _, keyItem := range keyTable {
		if !keyItem.Enable {
			continue
		}
		if keyItem.Name == name && keyItem.Key == key {
			return keyItem, nil
		}
	}
	return Key{}, &Error{"Permission denied"}
}

func getTTL(keyItemTTL int, defaultTTL int) int {
	if keyItemTTL != 0 {
		return keyItemTTL
	}
	return defaultTTL
}
