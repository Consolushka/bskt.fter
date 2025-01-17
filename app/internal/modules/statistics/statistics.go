package statistics

import (
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/internal/abstract"
	mlbl "IMP/app/internal/modules/statistics/internal/leagues/mlbl/repositories_factory"
	"IMP/app/internal/modules/statistics/internal/leagues/nba"
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
