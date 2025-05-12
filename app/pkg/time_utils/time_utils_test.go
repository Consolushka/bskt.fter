package time_utils

import (
	"testing"
)

// TestFormattedMinutesToSeconds tests correct formatting time in string formats to seconds
// Verify that pattern %m:%s w/o leading 0 is correctly converted to seconds
// Verify that pattern %m:%s with leading 0 is correctly converted to seconds
// Verify that pattern %m:%s with zero seconds is correctly converted to seconds
// Verify that patterns with text ('%m %s', '%m minutes, %s seconds') are correctly converted to seconds
// Verify that strings without seconds or minutes returns error with message "pattern must contain both %m and %s"
// Verify that invalid string format returns error with message "time string does not match the pattern"
// Verify that pattern with reversed minutes and seconds are correctly converted to seconds
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
			errorMessage: "time string does not match the pattern",
		},
		{
			name:         "invalid time string",
			timeStr:      "abc",
			pattern:      "%m:%s",
			errorMessage: "time string does not match the pattern",
		},
		{
			name:     "reversed pattern",
			timeStr:  "45:2",
			pattern:  "%s:%m",
			expected: 2*60 + 45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormattedMinutesToSeconds(tt.timeStr, tt.pattern)
			if tt.errorMessage != "" {
				if err == nil || err.Error() != tt.errorMessage {
					t.Errorf("FormattedMinutesToSeconds(%q, %q) error = %v, expected %v",
						tt.timeStr, tt.pattern, err, tt.errorMessage)
				}
				return
			}
			if err != nil {
				t.Errorf("FormattedMinutesToSeconds(%q, %q) error = %v, expected nil",
					tt.timeStr, tt.pattern, err)
				return
			}
			if result != tt.expected {
				t.Errorf("FormattedMinutesToSeconds(%q, %q) = %d, expected %d",
					tt.timeStr, tt.pattern, result, tt.expected)
			}
		})
	}
}
