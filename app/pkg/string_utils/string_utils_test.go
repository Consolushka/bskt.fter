package string_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetBoundaries verifies the behavior of the getBoundaries method
// for different language types:
// - Verify that when Language is set to Latin, the method returns the correct Unicode range (0-127) and no error
// - Verify that when Language is set to Cyrillic, the method returns the correct Unicode range (1024-1279) and no error
// - Verify that when Language is set to an invalid value, the method returns an error
func TestGetBoundaries(t *testing.T) {
	testCases := []struct {
		name        string
		language    Language
		expectedMin int32
		expectedMax int32
		expectError bool
	}{
		{
			name:        "Latin language",
			language:    Latin,
			expectedMin: 0,
			expectedMax: 127,
			expectError: false,
		},
		{
			name:        "Cyrillic language",
			language:    Cyrillic,
			expectedMin: 1024,
			expectedMax: 1279,
			expectError: false,
		},
		{
			name:        "Invalid language",
			language:    Language(999),
			expectedMin: 0,
			expectedMax: 0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minBoundaries, maxBoundaries, err := tc.language.getBoundaries()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedMin, minBoundaries)
			assert.Equal(t, tc.expectedMax, maxBoundaries)
		})
	}
}

// TestHasNonLanguageChars verifies the behavior of the HasNonLanguageChars function
// with different languages and inputs:
// - Verify that the function correctly identifies strings with and without non-Latin characters
// - Verify that the function correctly identifies strings with and without non-Cyrillic characters
// - Verify that when an invalid language is provided, the function returns an error
func TestHasNonLanguageChars(t *testing.T) {
	testCases := []struct {
		name           string
		text           string
		language       Language
		expectedResult bool
		expectError    bool
	}{
		// Latin language test cases
		{
			name:           "Only Latin characters",
			text:           "Hello World",
			language:       Latin,
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Contains non-Latin characters (Cyrillic)",
			text:           "Hello Привет",
			language:       Latin,
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Empty string with Latin",
			text:           "",
			language:       Latin,
			expectedResult: false,
			expectError:    false,
		},

		// Cyrillic language test cases
		{
			name:           "Only Cyrillic characters",
			text:           "привет мир",
			language:       Cyrillic,
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Contains non-Cyrillic characters (Latin)",
			text:           "Привет Hello",
			language:       Cyrillic,
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Empty string with Cyrillic",
			text:           "",
			language:       Cyrillic,
			expectedResult: false,
			expectError:    false,
		},

		// Invalid language test case
		{
			name:           "Invalid language",
			text:           "Hello",
			language:       Language(999),
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewStringUtils().HasNonLanguageChars(tc.text, tc.language)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

// TestRemovePunctuationAndSpaces verifies the behavior of the RemovePunctuationAndSpaces function
// with different inputs:
// - Verify that the function correctly removes spaces and punctuation marks from Latin text
// - Verify that the function correctly removes spaces and punctuation marks from Cyrillic text
// - Verify that the function correctly removes spaces and punctuation marks from mixed text
// - Verify that the function correctly handles strings containing only punctuation and spaces
// - Verify that the function correctly handles empty strings
// - Verify that the function correctly handles text with numbers and special characters
func TestRemovePunctuationAndSpaces(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "Latin text with spaces and punctuation",
			input:          "Hello, World! How are you?",
			expectedOutput: "HelloWorldHowareyou",
		},
		{
			name:           "Cyrillic text with spaces and punctuation",
			input:          "Привет, мир! Как дела?",
			expectedOutput: "ПриветмирКакдела",
		},
		{
			name:           "Mixed text with spaces and punctuation",
			input:          "Hello, Привет! How are you? Как дела?",
			expectedOutput: "HelloПриветHowareyouКакдела",
		},
		{
			name:           "Text with only punctuation and spaces",
			input:          "!@#$%^&*() ,.:;\"'[]{}",
			expectedOutput: "",
		},
		{
			name:           "Empty string",
			input:          "",
			expectedOutput: "",
		},
		{
			name:           "Text with numbers and special characters",
			input:          "Phone: +7(123)456-78-90, Email: test@example.com",
			expectedOutput: "Phone71234567890Emailtestexamplecom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewStringUtils().RemovePunctuationAndSpaces(tc.input)
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}
