package repositories_factory

import (
	"FTER/app/internal/statistics/abstract"
	"FTER/app/internal/statistics/leagues/nba/repositories/nba.com_api"
	"FTER/app/internal/statistics/leagues/nba/repositories/sport_radar"
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
		return nba_com_api.NewRepository()
	default:
		return nil
	}
}
