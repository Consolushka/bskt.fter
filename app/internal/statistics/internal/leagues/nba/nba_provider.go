package nba

import (
	"IMP/app/internal/modules/statistics/internal/leagues/nba/providers/sport_radar"
	"IMP/app/internal/statistics/internal/abstract"
	"IMP/app/internal/statistics/internal/leagues/nba/providers/cdn_nba"
	"os"
)

const (
	SPORTRADAR = "SPORTRADAR"
	NBA        = "CND.NBA"
)

// NewNbaStatsProvider based on provider returns provider for statistics
func NewNbaStatsProvider() abstract.StatsProvider {
	provider := os.Getenv("NBA_STATS_PROVIDER")

	switch provider {
	case SPORTRADAR:
		return sport_radar.NewProvider()
	case NBA:
		return cdn_nba.NewProvider()
	default:
		return nil
	}
}
