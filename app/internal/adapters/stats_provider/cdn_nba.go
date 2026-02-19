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
	return players.PlayerBioEntity{}, errors.New("CDN_NBA GetPlayerBio not implemented")
}

func (c CdnNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	return nil, errors.New("CDN_NBA GetGamesStatsByPeriod")
}

func (c CdnNbaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	return game, nil
}
