package stats_provider

import (
	"IMP/app/internal/core/games"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCdnNbaStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
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

func TestCdnNbaStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	adapter := CdnNbaStatsProviderAdapter{}
	bio, err := adapter.GetPlayerBio("any-id")

	assert.ErrorContains(t, err, "CDN_NBA GetPlayerBio not implemented")
	assert.Empty(t, bio.FullName)
}

func TestCdnNbaStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	adapter := CdnNbaStatsProviderAdapter{}
	game := games.GameStatEntity{GameModel: games.GameModel{Title: "TEST"}}

	result, err := adapter.EnrichGameStats(game)

	require.NoError(t, err)
	assert.Equal(t, game, result)
}
