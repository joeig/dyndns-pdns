package tools

// GetTTL selects the proper value of two choices
func GetTTL(keyItemTTL uint32, defaultTTL uint32) uint32 {
	if keyItemTTL != 0 {
		return keyItemTTL
	}
	return defaultTTL
}
