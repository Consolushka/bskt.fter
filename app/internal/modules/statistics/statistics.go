package statistics

import (
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/internal/abstract"
	"IMP/app/internal/modules/statistics/internal/leagues/mlbl"
	"IMP/app/internal/modules/statistics/internal/leagues/nba"
)

func NewLeagueProvider(league enums.League) abstract.StatsProvider {
	switch league {
	case enums.NBA:
		return nba.NewNbaStatsProvider()
	case enums.MLBL:
		return mlbl.NewMLBLProvider()
	default:
		return nil
	}
}
