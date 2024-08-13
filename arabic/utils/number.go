package utils

// Min returns the smaller of two integer values.
// This utility function is frequently used in comparing indices or lengths.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two integer values.
// This utility function is frequently used in comparing indices or lengths.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
