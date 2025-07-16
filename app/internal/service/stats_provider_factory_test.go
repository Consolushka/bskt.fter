package service

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/ports"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewStatsProvider tests the NewStatsProvider function
// Verify that when NBA league name is provided while creating stats provider - returns CdnNbaStatsProviderAdapter and no error
// Verify that when MLBL league name is provided while creating stats provider - returns InfobasketStatsProviderAdapter and no error
// Verify that when unknown league name is provided while creating stats provider - returns nil and specific error message
func TestNewStatsProvider(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		expected ports.StatsProvider
		errorMsg string
	}{
		{
			name:     "successfully creates NBA stats provider",
			data:     "NBA",
			expected: stats_provider.CdnNbaStatsProviderAdapter{},
			errorMsg: "",
		},
		{
			name:     "successfully creates MLBL stats provider",
			data:     "MLBL",
			expected: stats_provider.InfobasketStatsProviderAdapter{},
			errorMsg: "",
		},
		{
			name:     "returns error for unknown league name - differs from supported leagues",
			data:     "UNKNOWN",
			expected: nil,
			errorMsg: "unknown league: UNKNOWN",
		},
		{
			name:     "returns error for empty league name - differs from supported leagues",
			data:     "",
			expected: nil,
			errorMsg: "unknown league: ",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := NewStatsProvider(tc.data)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
