package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/ports"
)

type MlblStatsProviderFactory struct {
}

func (m MlblStatsProviderFactory) ProviderName() string {
	return "INFOBASKET"
}

func (m MlblStatsProviderFactory) Create() (ports.StatsProvider, error) {
	return stats_provider.InfobasketStatsProviderAdapter{}, nil
}
