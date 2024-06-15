package fshare

// RemoveDuplicates removes duplicate values from a map
func (s *Fshare) RemoveDuplicates(inputMap map[string]*QueryDetailResponse) map[string]*QueryDetailResponse {
	uniqueMap := make(map[string]*QueryDetailResponse)

	// Iterate through the original map and add unique key-value pairs to the new map
	for key, value := range inputMap {
		// Check if the key is already present in the uniqueMap
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = value
		}
	}

	return uniqueMap
}
