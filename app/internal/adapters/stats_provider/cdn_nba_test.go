package stats_provider

import (
	"IMP/app/internal/core/games"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCdnNbaStatsProviderAdapter_GetGamesStatsByDate tests the GetGamesStatsByPeriod method
// Verify that when any date is provided while calling GetGamesStatsByPeriod - returns nil and specific error message
func TestCdnNbaStatsProviderAdapter_GetGamesStatsByDate(t *testing.T) {
	cases := []struct {
		name     string
		data     time.Time
		expected []games.GameStatEntity
		errorMsg string
	}{
		{
			name:     "returns error for any date input",
			data:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: nil,
			errorMsg: "CDN_NBA GetGamesStatsByPeriod",
		},
	}

	adapter := CdnNbaStatsProviderAdapter{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := adapter.GetGamesStatsByPeriod(tc.data, tc.data)

			if tc.errorMsg != "" {
				require.EqualError(t, err, tc.errorMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
