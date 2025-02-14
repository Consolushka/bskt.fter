package statistics

import (
	"IMP/app/internal/modules/statistics/internal/abstract"
	"IMP/app/internal/modules/statistics/internal/leagues/mlbl"
	"IMP/app/internal/modules/statistics/internal/leagues/nba"
)

func NewLeagueProvider(leagueAliasEn string) abstract.StatsProvider {
	switch leagueAliasEn {
	case "NBA":
		return nba.NewNbaStatsProvider()
	case "MLBL":
		return mlbl.NewMLBLProvider()
	default:
		return nil
	}
}
