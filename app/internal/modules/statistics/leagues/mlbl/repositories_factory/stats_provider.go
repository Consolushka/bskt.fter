package repositories_factory

import (
	"IMP/app/internal/modules/statistics/abstract"
	"IMP/app/internal/modules/statistics/leagues/mlbl/repositories/infobasket"
)

const (
	INFOBASKET = "INFOBASKET"
)

// NewMLBLRepository based on provider returns repository for statistics from MLBL league
func NewMLBLRepository() abstract.StatsRepository {
	provider := "INFOBASKET"

	switch provider {
	case INFOBASKET:
		return infobasket.NewRepository()
	default:
		return nil
	}
}
