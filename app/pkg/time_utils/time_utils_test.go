package time_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestFormattedMinutesToSeconds verifies the behavior of the FormattedMinutesToSeconds method
// under various input conditions:
// - Verify that pattern %m:%s without leading zeros is correctly converted to seconds
// - Verify that pattern %m:%s with leading zeros is correctly converted to seconds
// - Verify that pattern %m:%s with zero seconds is correctly converted to seconds
// - Verify that pattern %m:%s with zero minutes is correctly converted to seconds
// - Verify that patterns with additional text ('%mm %ss') are correctly converted to seconds
// - Verify that patterns with descriptive text ('Time: %mm %ss remaining') are correctly converted to seconds
// - Verify that strings without minutes or seconds returns error with message "pattern must contain both %m and %s"
// - Verify that invalid string format returns error with message "time string does not match the pattern"
// - Verify that patterns with reversed minutes and seconds are correctly converted to seconds
// - Verify that complex patterns like ISO-8601 duration format are correctly parsed
func TestFormattedMinutesToSeconds(t *testing.T) {
	tests := []struct {
		name         string
		timeStr      string
		pattern      string
		expected     int
		errorMessage string
	}{
		{
			name:     "basic minutes and seconds",
			timeStr:  "3:45",
			pattern:  "%m:%s",
			expected: 3*60 + 45,
		},
		{
			name:     "single digit minutes and seconds",
			timeStr:  "1:05",
			pattern:  "%m:%s",
			expected: 1*60 + 5,
		},
		{
			name:     "zero minutes",
			timeStr:  "0:30",
			pattern:  "%m:%s",
			expected: 30,
		},
		{
			name:     "zero seconds",
			timeStr:  "5:00",
			pattern:  "%m:%s",
			expected: 5 * 60,
		},
		{
			name:     "different pattern format",
			timeStr:  "10m 15s",
			pattern:  "%mm %ss",
			expected: 10*60 + 15,
		},
		{
			name:     "pattern with additional text",
			timeStr:  "Time: 2m 30s remaining",
			pattern:  "Time: %mm %ss remaining",
			expected: 2*60 + 30,
		},
		{
			name:         "missing minutes in pattern",
			timeStr:      "45s",
			pattern:      "%ss",
			errorMessage: "pattern must contain both %m and %s",
		},
		{
			name:         "missing seconds in pattern",
			timeStr:      "5m",
			pattern:      "%mm",
			errorMessage: "pattern must contain both %m and %s",
		},
		{
			name:         "timeStr shorter than pattern",
			timeStr:      "1:2",
			pattern:      "%m:%s extra",
			errorMessage: "time string '1:2' does not match the pattern '%m:%s extra'",
		},
		{
			name:         "invalid time string",
			timeStr:      "abc",
			pattern:      "%m:%s",
			errorMessage: "time string 'abc' does not match the pattern '%m:%s'",
		},
		{
			name:     "reversed pattern",
			timeStr:  "45:2",
			pattern:  "%s:%m",
			expected: 2*60 + 45,
		},
		{
			name:     "reversed pattern",
			timeStr:  "PT22M28.44S",
			pattern:  "PT%mM%sS",
			expected: 22*60 + 28,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTimeUtils().FormattedMinutesToSeconds(tt.timeStr, tt.pattern)

			if tt.errorMessage != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
