package factory

import (
	"FTER/app/internal/modules/statistics/abstract"
	"FTER/app/internal/modules/statistics/enums"
	mlbl "FTER/app/internal/modules/statistics/leagues/mlbl/repositories_factory"
	nba "FTER/app/internal/modules/statistics/leagues/nba/repositories_factory"
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
