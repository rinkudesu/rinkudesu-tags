package utils

func Contains[T comparable](haystack []T, needle T) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}
