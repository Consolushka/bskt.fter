package imp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPERs_Order verifies the behavior of the Order method
// in the PERs type under various conditions:
// - Verify that Clean PER returns order value 0
// - Verify that Bench PER returns order value 1
// - Verify that Start PER returns order value 2
// - Verify that FullGame PER returns order value 3
// - Verify that invalid PER returns order value -1
func TestPERs_Order(t *testing.T) {
	cases := []struct {
		name     string
		per      PERs
		expected int
	}{
		{
			name:     "Clean PER Order",
			per:      Clean,
			expected: 0,
		},
		{
			name:     "Bench PER Order",
			per:      Bench,
			expected: 1,
		},
		{
			name:     "Start PER Order",
			per:      Start,
			expected: 2,
		},
		{
			name:     "FullGame PER Order",
			per:      FullGame,
			expected: 3,
		},
		{
			name:     "Invalid PER Order",
			per:      "Unexpected",
			expected: -1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.per.Order()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// TestImp_EvaluateClean verifies the behavior of the EvaluateClean function
// under various input conditions:
// - Verify calculation with positive values for both plsMin and finalDifference
// - Verify calculation with negative plsMin and positive finalDifference
// - Verify calculation with both negative plsMin and finalDifference
// - Verify that zero playedSeconds returns zero regardless of other parameters
func TestImp_EvaluateClean(t *testing.T) {
	cases := []struct {
		name                  string
		playedSeconds         int
		plsMin                int
		finalDifference       int
		fullGameTimeInMinutes int
		expected              float64
	}{
		{
			name:                  "Evaluate clean IMP with positive both plsMin and finalDifference",
			playedSeconds:         1800,
			plsMin:                7,
			finalDifference:       10,
			fullGameTimeInMinutes: 48,
			expected:              0.024999999999999994,
		},
		{
			name:                  "Evaluate clean IMP with negative plsMin and positive finalDifference",
			playedSeconds:         1800,
			plsMin:                -7,
			finalDifference:       10,
			fullGameTimeInMinutes: 48,
			expected:              -0.44166666666666665,
		},
		{
			name:                  "Evaluate clean IMP with both negative plsMin and finalDifference",
			playedSeconds:         1800,
			plsMin:                -7,
			finalDifference:       -15,
			fullGameTimeInMinutes: 48,
			expected:              0.07916666666666666,
		},
		{
			name:                  "Evaluate clean IMP with zero playedSeconds",
			playedSeconds:         0,
			plsMin:                0,
			finalDifference:       10,
			fullGameTimeInMinutes: 48,
			expected:              0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := EvaluateClean(tc.playedSeconds, tc.plsMin, tc.finalDifference, tc.fullGameTimeInMinutes)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// TestPERs_ToString verifies the behavior of the ToString method
// in the PERs type for all valid PER values:
// - Verify that Clean PER returns "Clean" string
// - Verify that Bench PER returns "Bench" string
// - Verify that Start PER returns "Start" string
// - Verify that FullGame PER returns "FullGame" string
func TestPERs_ToString(t *testing.T) {
	cases := []struct {
		name     string
		per      PERs
		expected string
	}{
		{
			name:     "Clean PER ToString",
			per:      Clean,
			expected: "Clean",
		},
		{
			name:     "Bench PER ToString",
			per:      Bench,
			expected: "Bench",
		},
		{
			name:     "Start PER ToString",
			per:      Start,
			expected: "Start",
		},
		{
			name:     "FullGame PER ToString",
			per:      FullGame,
			expected: "FullGame",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.per.ToString()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
