package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/infra/sportoteka"
	"IMP/app/internal/ports"
)

type UbaStatsProviderFactory struct {
	Tag  string
	Year int
}

func (m UbaStatsProviderFactory) ProviderName() string {
	return "SPORTOTEKA"
}

func (m UbaStatsProviderFactory) Create() (ports.StatsProvider, error) {
	return stats_provider.NewSportotekaStatsProvider(
		sportoteka.NewClient(),
		m.Tag,
		m.Year,
	), nil
}
