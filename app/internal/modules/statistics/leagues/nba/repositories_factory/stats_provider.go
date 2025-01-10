package repositories_factory

import (
	"IMP/app/internal/modules/statistics/abstract"
	"IMP/app/internal/modules/statistics/leagues/nba/repositories/cdn_nba"
	"IMP/app/internal/modules/statistics/leagues/nba/repositories/sport_radar"
)

const (
	SPORTRADAR = "SPORTRADAR"
	NBA        = "NBA"
)

// NewNbaStatsRepository based on provider returns repository for statistics
func NewNbaStatsRepository() abstract.StatsRepository {
	provider := "NBA"

	switch provider {
	case SPORTRADAR:
		return sport_radar.NewRepository()
	case NBA:
		return cdn_nba.NewRepository()
	default:
		return nil
	}
}
