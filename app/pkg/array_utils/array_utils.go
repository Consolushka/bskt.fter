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

func Map[T any, R any](array []T, callback func(T) (R, error)) ([]R, error) {
	mapped := make([]R, 0)
	for _, item := range array {
		result, err := callback(item)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, result)
	}
	return mapped, nil
}
