package tools

// GetTTL selects the proper value of two choices
func GetTTL(keyItemTTL int, defaultTTL int) int {
	if keyItemTTL != 0 {
		return keyItemTTL
	}
	return defaultTTL
}
