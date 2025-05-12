package array_utils_test

import (
	"IMP/app/pkg/array_utils"
	"reflect"
	"testing"
)

// TestFilterInts tests the Filter function with integer slices.
// Verify that when filtering a slice of integers for even numbers, only even numbers are returned.
// Verify that when filtering an empty slice, an empty slice is returned.
func TestFilterInts(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		filter   func(int) bool
		expected []int
	}{
		{
			name:     "filter even numbers",
			input:    []int{1, 2, 3, 4, 5, 6},
			filter:   func(n int) bool { return n%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "filter numbers greater than 3",
			input:    []int{1, 2, 3, 4, 5, 6},
			filter:   func(n int) bool { return n > 3 },
			expected: []int{4, 5, 6},
		},
		{
			name:     "filter with no matches",
			input:    []int{1, 2, 3, 4},
			filter:   func(n int) bool { return n > 10 },
			expected: []int{},
		},
		{
			name:     "filter empty slice",
			input:    []int{},
			filter:   func(n int) bool { return n > 0 },
			expected: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := array_utils.Filter(tc.input, tc.filter)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Filter() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestFilterStrings tests the Filter function with string slices.
// Verify that when filtering a slice of strings based on length, only strings meeting the length criteria are returned.
func TestFilterStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		filter   func(string) bool
		expected []string
	}{
		{
			name:     "filter strings longer than 3 characters",
			input:    []string{"a", "ab", "abc", "abcd", "abcde"},
			filter:   func(s string) bool { return len(s) > 3 },
			expected: []string{"abcd", "abcde"},
		},
		{
			name:     "filter strings containing 'a'",
			input:    []string{"apple", "banana", "orange", "grape", "kiwi"},
			filter:   func(s string) bool { return len(s) > 0 && s[0] == 'a' },
			expected: []string{"apple"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := array_utils.Filter(tc.input, tc.filter)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Filter() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestFilterStructs tests the Filter function with struct slices.
// Verify that when filtering a slice of structs based on a field value, only structs meeting the criteria are returned.
func TestFilterStructs(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	people := []person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Dave", Age: 40},
	}

	tests := []struct {
		name     string
		input    []person
		filter   func(person) bool
		expected []person
	}{
		{
			name:  "filter people older than 30",
			input: people,
			filter: func(p person) bool {
				return p.Age > 30
			},
			expected: []person{
				{Name: "Charlie", Age: 35},
				{Name: "Dave", Age: 40},
			},
		},
		{
			name:  "filter people with names starting with 'A'",
			input: people,
			filter: func(p person) bool {
				return len(p.Name) > 0 && p.Name[0] == 'A'
			},
			expected: []person{
				{Name: "Alice", Age: 25},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := array_utils.Filter(tc.input, tc.filter)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Filter() = %v, want %v", result, tc.expected)
			}
		})
	}
}
