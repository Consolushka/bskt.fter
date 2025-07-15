package stats_provider

import (
	"IMP/app/internal/core/games"
	"errors"
	"time"
)

type InfobasketStatsProviderAdapter struct {
}

func (i InfobasketStatsProviderAdapter) GetGamesStatsByDate(date time.Time) ([]games.GameStatEntity, error) {
	return nil, errors.New("CDN_NBA infobasket")
}
