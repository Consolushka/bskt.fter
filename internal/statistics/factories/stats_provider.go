package factories

import (
	"FTER/internal/statistics/repositories"
	"FTER/internal/statistics/repositories/sport_radar/repository"
	"errors"
)

const (
	SPORTRADAR = "SPORTRADAR"
)

// NewStatsRepository based on provider returns repository for statistics
func NewStatsRepository() (repositories.StatsRepository, error) {
	provider := "SPORTRADAR"

	switch provider {
	case SPORTRADAR:
		return repository.NewSportRadarRepository(), nil
	default:
		return nil, errors.New("unknown stats provider")
	}
}
