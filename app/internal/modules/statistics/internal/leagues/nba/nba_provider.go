package nba

import (
	"IMP/app/internal/modules/statistics/internal/abstract"
	"IMP/app/internal/modules/statistics/internal/leagues/nba/providers/cdn_nba"
	"IMP/app/internal/modules/statistics/internal/leagues/nba/providers/sport_radar"
)

const (
	SPORTRADAR = "SPORTRADAR"
	NBA        = "NBA"
)

// NewNbaStatsProvider based on provider returns provider for statistics
func NewNbaStatsProvider() abstract.StatsProvider {
	provider := "NBA"

	switch provider {
	case SPORTRADAR:
		return sport_radar.NewProvider()
	case NBA:
		return cdn_nba.NewProvider()
	default:
		return nil
	}
}
