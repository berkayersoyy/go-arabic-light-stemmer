package utils

// MinFromSlice finds and returns the minimum value from a slice of integers.
// This utility function is commonly used in determining the smallest index or position.
func MinFromSlice(ints []int) int {
	minVal := ints[0]
	for _, val := range ints {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// MaxFromSlice finds and returns the maximum value from a slice of integers.
// This utility function is commonly used in determining the largest index or position.
func MaxFromSlice(ints []int) int {
	maxVal := ints[0]
	for _, val := range ints {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// AffixInList checks if a given affix is present in a predefined list of valid affixes.
// This helper function is essential for validating whether a prefix-suffix combination is acceptable.
func AffixInList(affix string, list []string) bool {
	for _, a := range list {
		if a == affix {
			return true
		}
	}
	return false
}

// Contains checks if a slice contains a specific string item.
// This utility function is useful for validating membership in lists or sets.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
