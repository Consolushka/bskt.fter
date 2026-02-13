package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"errors"
	"time"
)

type CdnNbaStatsProviderAdapter struct {
}

func (c CdnNbaStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (c CdnNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	return nil, errors.New("CDN_NBA GetGamesStatsByPeriod")
}

func (c CdnNbaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	return game, nil
}
