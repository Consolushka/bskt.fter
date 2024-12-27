package factory

import (
	"IMP/app/internal/modules/statistics/abstract"
	"IMP/app/internal/modules/statistics/enums"
	mlbl "IMP/app/internal/modules/statistics/leagues/mlbl/repositories_factory"
	nba "IMP/app/internal/modules/statistics/leagues/nba/repositories_factory"
)

func NewLeagueRepository(league enums.League) abstract.StatsRepository {
	switch league {
	case enums.NBA:
		return nba.NewNbaStatsRepository()
	case enums.MLBL:
		return mlbl.NewMLBLRepository()
	default:
		return nil
	}
}
