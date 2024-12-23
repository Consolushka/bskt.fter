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

func Map[T any, R any](array []T, callback func(T) R) []R {
	mapped := make([]R, 0)
	for _, item := range array {
		mapped = append(mapped, callback(item))
	}
	return mapped
}
