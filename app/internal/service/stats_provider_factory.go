package service

import (
	"IMP/app/internal/adapters/stats_provider_factory"
	"IMP/app/internal/ports"
	"errors"
)

func NewStatsProvider(leagueName string) (ports.StatsProvider, error) {
	var factory ports.StatsProviderFactory

	switch leagueName {
	case "NBA":
		factory = stats_provider_factory.NbaStatsProviderFactory{}
	case "MLBL":
		factory = stats_provider_factory.MlblStatsProviderFactory{}
	default:
		return nil, errors.New("unknown league: " + leagueName)
	}

	return factory.Create()
}
