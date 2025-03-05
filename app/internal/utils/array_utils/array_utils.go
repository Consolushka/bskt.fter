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

// Sum calculates the sum of values returned by the accessor function for each item in the array.
// The accessor function extracts a numeric value from an item of type T.
func Sum[T any, N interface {
	int | int64 | float64 | float32
}](array []T, accessor func(T) N) N {
	var sum N
	for _, item := range array {
		sum += accessor(item)
	}
	return sum
}

// Average calculates the average of values returned by the accessor function for each item in the array.
// Returns 0 if the array is empty.
func Average[T any, N interface {
	int | int64 | float64 | float32
}](array []T, accessor func(T) N) float64 {
	if len(array) == 0 {
		return 0
	}

	var sum N
	for _, item := range array {
		sum += accessor(item)
	}

	return float64(sum) / float64(len(array))
}
