package factories

import (
	"FTER/internal/statistics/repositories"
	nbaRepository "FTER/internal/statistics/repositories/nba/repository"
	sportRadarRepository "FTER/internal/statistics/repositories/sport_radar/repository"
	"errors"
)

const (
	SPORTRADAR = "SPORTRADAR"
	NBA        = "NBA"
)

// NewStatsRepository based on provider returns repository for statistics
func NewStatsRepository() (repositories.StatsRepository, error) {
	provider := "NBA"

	switch provider {
	case SPORTRADAR:
		return sportRadarRepository.NewSportRadarRepository(), nil
	case NBA:
		return nbaRepository.NewNbaRepository(), nil
	default:
		return nil, errors.New("unknown stats provider")
	}
}
