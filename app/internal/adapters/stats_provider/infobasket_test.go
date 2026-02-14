package stats_provider

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestInfobasketStatsProviderAdapter_GetGamesStatsByDate tests the GetGamesStatsByPeriod method
// Verify that when any date is provided while calling GetGamesStatsByPeriod - returns nil and specific error message
func TestInfobasketStatsProviderAdapter_GetGamesStatsByDate(t *testing.T) {
	adapter := InfobasketStatsProviderAdapter{}
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	assert.Panics(t, func() {
		_, _ = adapter.GetGamesStatsByPeriod(date, date)
	})
}
