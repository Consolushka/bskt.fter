package arrays

func Filter[T any](array []T, callback func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range array {
		if callback(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
