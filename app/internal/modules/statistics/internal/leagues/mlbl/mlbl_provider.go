package mlbl

import (
	"IMP/app/internal/modules/statistics/internal/abstract"
	"IMP/app/internal/modules/statistics/internal/leagues/mlbl/providers/infobasket"
)

const (
	INFOBASKET = "INFOBASKET"
)

// NewMLBLProvider based on provider returns provider for statistics from MLBL league
func NewMLBLProvider() abstract.StatsProvider {
	provider := "INFOBASKET"

	switch provider {
	case INFOBASKET:
		return infobasket.NewProvider()
	default:
		return nil
	}
}
