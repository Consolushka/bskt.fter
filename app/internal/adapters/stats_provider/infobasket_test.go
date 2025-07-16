package stats_provider

import (
	"IMP/app/internal/core/games"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestInfobasketStatsProviderAdapter_GetGamesStatsByDate tests the GetGamesStatsByDate method
// Verify that when any date is provided while calling GetGamesStatsByDate - returns nil and specific error message
func TestInfobasketStatsProviderAdapter_GetGamesStatsByDate(t *testing.T) {
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
			errorMsg: "CDN_NBA infobasket",
		},
	}

	adapter := InfobasketStatsProviderAdapter{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := adapter.GetGamesStatsByDate(tc.data)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
