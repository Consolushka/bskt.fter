package array_utils

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

func Sort[T any](array []T, less func(a, b T) bool) []T {
	sorted := make([]T, len(array))
	copy(sorted, array)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if !less(sorted[j], sorted[j+1]) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}
