package statistics

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/statistics/internal/abstract"
	"IMP/app/internal/statistics/internal/leagues/mlbl"
	"IMP/app/internal/statistics/internal/leagues/nba"
	"strings"
)

func NewLeagueProvider(leagueAliasEn string) abstract.StatsProvider {
	switch leagueAliasEn {
	case strings.ToUpper(domain.NBAAlias):
		return nba.NewNbaStatsProvider()
	case strings.ToUpper(domain.MLBLAlias):
		return mlbl.NewProvider()
	default:
		return nil
	}
}
