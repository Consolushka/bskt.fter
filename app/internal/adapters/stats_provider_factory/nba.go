package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/ports"
)

type NbaStatsProviderFactory struct {
}

func (m NbaStatsProviderFactory) Create() (ports.StatsProvider, error) {
	return stats_provider.CdnNbaStatsProviderAdapter{}, nil
}
