package array_utils_test

import (
	"IMP/app/pkg/array_utils"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// TestFilterInts verifies the behavior of the Filter function with integer slices
// under various conditions:
// - Verify that when filtering a slice of integers for even numbers, only even numbers are returned
// - Verify that when filtering a slice of integers for numbers greater than a threshold, only matching numbers are returned
// - Verify that when filtering with criteria that match no elements, an empty slice is returned
// - Verify that when filtering an empty slice, an empty slice is returned
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
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestFilterStrings verifies the behavior of the Filter function with string slices
// under various conditions:
// - Verify that when filtering a slice of strings based on length, only strings meeting the length criteria are returned
// - Verify that when filtering a slice of strings based on starting character, only strings with matching first character are returned
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
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestFilterStructs verifies the behavior of the Filter function with struct slices
// under various conditions:
// - Verify that when filtering a slice of structs based on a numeric field value, only structs meeting the criteria are returned
// - Verify that when filtering a slice of structs based on a string field value, only structs meeting the criteria are returned
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
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMapInts verifies the behavior of the Map function with integer slices
// under various conditions:
// - Verify that when mapping a slice of integers with a doubling transformation, each element is correctly doubled
// - Verify that when mapping a slice of integers with an addition transformation, each element is correctly increased
// - Verify that when mapping an empty slice, an empty slice is returned
// - Verify that when mapping encounters an error during transformation, the error is properly returned
func TestMapInts(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapper   func(int) (int, error)
		expected []int
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "double each number",
			input:    []int{1, 2, 3, 4, 5},
			mapper:   func(n int) (int, error) { return n * 2, nil },
			expected: []int{2, 4, 6, 8, 10},
			wantErr:  false,
		},
		{
			name:     "add 10 to each number",
			input:    []int{1, 2, 3, 4, 5},
			mapper:   func(n int) (int, error) { return n + 10, nil },
			expected: []int{11, 12, 13, 14, 15},
			wantErr:  false,
		},
		{
			name:     "map empty slice",
			input:    []int{},
			mapper:   func(n int) (int, error) { return n * 2, nil },
			expected: []int{},
			wantErr:  false,
		},
		{
			name:  "error on negative number",
			input: []int{1, 2, -3, 4, 5},
			mapper: func(n int) (int, error) {
				if n < 0 {
					return 0, errors.New("negative number not allowed")
				}
				return n * 2, nil
			},
			expected: nil,
			wantErr:  true,
			errMsg:   "negative number not allowed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := array_utils.Map(tc.input, tc.mapper)

			if tc.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// TestMapIntToString verifies the behavior of the Map function with type conversion
// from integers to strings:
// - Verify that when mapping a slice of integers to strings, each integer is correctly converted to its string representation
func TestMapIntToString(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	result, err := array_utils.Map(input, func(n int) (string, error) {
		return strconv.Itoa(n), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestMapStringToInt verifies the behavior of the Map function with type conversion
// from strings to integers:
// - Verify that when mapping a slice of strings to integers, each string is correctly parsed to its integer value
func TestMapStringToInt(t *testing.T) {
	input := []string{"1", "2", "3"}
	expected := []int{1, 2, 3}

	result, err := array_utils.Map(input, func(s string) (int, error) {
		return strconv.Atoi(s)
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestMapStringToIntError verifies the error handling of the Map function
// when string to integer conversion fails:
// - Verify that when mapping a slice of strings to integers encounters an invalid number, the appropriate error is returned
func TestMapStringToIntError(t *testing.T) {
	input := []string{"1", "2", "not-a-number", "3"}

	result, err := array_utils.Map(input, func(s string) (int, error) {
		return strconv.Atoi(s)
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMapStructs verifies the behavior of the Map function with struct slices
// under various conditions:
// - Verify that when mapping a slice of structs to another type, the transformation is correctly applied to each element
// - Verify that when mapping with validation that passes for all elements, all transformations succeed
// - Verify that when validation fails during mapping, the appropriate error is returned
func TestMapStructs(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	type personSummary struct {
		FullName string
		IsAdult  bool
	}

	people := []person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 35},
	}

	invalidPeople := []person{
		{Name: "Alice", Age: 25},
		{Name: "", Age: 30}, // Invalid name
		{Name: "Charlie", Age: 35},
	}

	tests := []struct {
		name     string
		input    []person
		mapper   func(person) (personSummary, error)
		expected []personSummary
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "successful struct transformation",
			input: people,
			mapper: func(p person) (personSummary, error) {
				return personSummary{
					FullName: p.Name,
					IsAdult:  p.Age >= 18,
				}, nil
			},
			expected: []personSummary{
				{FullName: "Alice", IsAdult: true},
				{FullName: "Bob", IsAdult: false},
				{FullName: "Charlie", IsAdult: true},
			},
			wantErr: false,
		},
		{
			name:  "struct transformation with validation - no errors",
			input: people,
			mapper: func(p person) (personSummary, error) {
				if p.Name == "" {
					return personSummary{}, errors.New("name cannot be empty")
				}
				if p.Age < 0 {
					return personSummary{}, fmt.Errorf("invalid age: %d", p.Age)
				}
				return personSummary{
					FullName: p.Name,
					IsAdult:  p.Age >= 18,
				}, nil
			},
			expected: []personSummary{
				{FullName: "Alice", IsAdult: true},
				{FullName: "Bob", IsAdult: false},
				{FullName: "Charlie", IsAdult: true},
			},
			wantErr: false,
		},
		{
			name:  "struct transformation with validation error",
			input: invalidPeople,
			mapper: func(p person) (personSummary, error) {
				if p.Name == "" {
					return personSummary{}, errors.New("name cannot be empty")
				}
				return personSummary{
					FullName: p.Name,
					IsAdult:  p.Age >= 18,
				}, nil
			},
			expected: nil,
			wantErr:  true,
			errMsg:   "name cannot be empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := array_utils.Map(tc.input, tc.mapper)

			if tc.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
