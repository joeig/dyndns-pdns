package pdns

func getTTL(keyItemTTL int, defaultTTL int) int {
	if keyItemTTL != 0 {
		return keyItemTTL
	}
	return defaultTTL
}
