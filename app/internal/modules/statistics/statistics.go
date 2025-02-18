package statistics

import (
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	"IMP/app/internal/modules/statistics/internal/abstract"
	"IMP/app/internal/modules/statistics/internal/leagues/mlbl"
	"IMP/app/internal/modules/statistics/internal/leagues/nba"
	"strings"
)

func NewLeagueProvider(leagueAliasEn string) abstract.StatsProvider {
	switch leagueAliasEn {
	case strings.ToUpper(leaguesModels.NBAAlias):
		return nba.NewNbaStatsProvider()
	case strings.ToUpper(leaguesModels.MLBLAlias):
		return mlbl.NewMLBLProvider()
	default:
		return nil
	}
}
