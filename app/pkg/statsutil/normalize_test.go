package statsutil

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPercentage100(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float32
	}{
		{"Normal 50.5", 50.5, 0.505},
		{"Zero", 0, 0},
		{"100%", 100, 1.0},
		{"Below Zero", -1, 0},
		{"Above 100", 101, 1.0},
		{"NaN", math.NaN(), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FromPercentage100(tt.input)
			assert.InDelta(t, tt.expected, res, 0.0001)
		})
	}
}

func TestFromRatio(t *testing.T) {
	tests := []struct {
		name      string
		made      int
		attempted int
		expected  float32
	}{
		{"Normal 1/2", 1, 2, 0.5},
		{"Normal 2/3", 2, 3, 0.6666667},
		{"Zero attempted", 1, 0, 0},
		{"Zero made", 0, 5, 0},
		{"Full", 5, 5, 1.0},
		{"Impossible more made than attempted", 6, 5, 1.0},
		{"Negative made", -1, 5, 0},
		{"Negative attempted", 1, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FromRatio(tt.made, tt.attempted)
			assert.InDelta(t, tt.expected, res, 0.0001)
		})
	}
}
