package repositories_factory

import (
	"FTER/app/internal/statistics/abstract"
	"FTER/app/internal/statistics/leagues/mlbl/repositories/infobasket"
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
