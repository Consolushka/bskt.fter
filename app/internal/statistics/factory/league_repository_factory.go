package factory

import (
	"FTER/app/internal/enums"
	"FTER/app/internal/statistics/abstract"
	mlbl "FTER/app/internal/statistics/leagues/mlbl/repositories_factory"
	nba "FTER/app/internal/statistics/leagues/nba/repositories_factory"
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
