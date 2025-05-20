package translator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestTranslator_Translate verifies the behavior of the Translate method
// in the translator struct under various conditions:
// - Verify that when source language is nil, translation works correctly
// - Verify that when source language is provided, translation works correctly
// - Verify that translation works for different language pairs
// - Verify that translation preserves special characters and formatting
func TestTranslator_Translate(t *testing.T) {
	cases := []struct {
		name       string
		source     string
		sourceLang *string
		targetLang string
		expected   string
	}{
		{
			name:       "Translate English to Russian with auto-detected source",
			source:     "Hello",
			sourceLang: nil,
			targetLang: "ru",
			expected:   "Привет",
		},
		{
			name:       "Translate English to Russian with specified source",
			source:     "Hello",
			sourceLang: stringPtr("en"),
			targetLang: "ru",
			expected:   "Привет",
		},
		{
			name:       "Translate Russian to English",
			source:     "Привет",
			sourceLang: stringPtr("ru"),
			targetLang: "en",
			expected:   "Hello",
		},
		{
			name:       "Translate English to Spanish",
			source:     "Hello",
			sourceLang: stringPtr("en"),
			targetLang: "es",
			expected:   "Hola",
		},
	}

	translator := NewTranslator()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := translator.Translate(tc.source, tc.sourceLang, tc.targetLang)

			// Since translations might vary slightly, we check if the expected text
			// is contained in the result or if the result is contained in the expected
			// This makes the test more robust against minor translation differences
			containsExpected := assert.Contains(t, result, tc.expected) ||
				assert.Contains(t, tc.expected, result)

			if !containsExpected {
				t.Errorf("Expected translation to contain '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

// TestTranslator_TranslatePreservesFormatting verifies that the translator
// preserves special characters and formatting in the translated text
func TestTranslator_TranslatePreservesFormatting(t *testing.T) {
	cases := []struct {
		name       string
		source     string
		sourceLang *string
		targetLang string
		preserved  []string
	}{
		{
			name:       "Preserve numbers in translation",
			source:     "I have 5 apples and 3 oranges",
			sourceLang: stringPtr("en"),
			targetLang: "ru",
			preserved:  []string{"5", "3"},
		},
		{
			name:       "Preserve special characters",
			source:     "Email: test@example.com, Phone: +1-234-567-8900",
			sourceLang: stringPtr("en"),
			targetLang: "ru",
			preserved:  []string{"test@example.com", "+1-234-567-8900"},
		},
	}

	translator := NewTranslator()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := translator.Translate(tc.source, tc.sourceLang, tc.targetLang)

			for _, preserved := range tc.preserved {
				assert.Contains(t, result, preserved,
					"Translation should preserve '%s' but got '%s'", preserved, result)
			}
		})
	}
}

// TestTranslator_TranslateEmptyString verifies that the translator
// handles empty strings correctly
func TestTranslator_TranslateEmptyString(t *testing.T) {
	translator := NewTranslator()
	result := translator.Translate("", stringPtr("en"), "ru")
	assert.Empty(t, result, "Translation of empty string should be empty")
}

// TestNewTranslator verifies that NewTranslator returns a valid translator instance
func TestNewTranslator(t *testing.T) {
	translator := NewTranslator()
	assert.NotNil(t, translator, "NewTranslator should return a non-nil translator")
	assert.Implements(t, (*Interface)(nil), translator, "NewTranslator should return an implementation of Interface")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
