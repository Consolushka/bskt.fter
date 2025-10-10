package stats_provider

import (
	"IMP/app/internal/core/games"
	"errors"
	"time"
)

type CdnNbaStatsProviderAdapter struct {
}

func (c CdnNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	return nil, errors.New("CDN_NBA GetGamesStatsByPeriod")
}
