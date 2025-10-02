package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/infra/infobasket"
	"IMP/app/internal/ports"
)

type FonbetSuperleagueStatsProviderFactory struct {
	ExternalId int
}

func (m FonbetSuperleagueStatsProviderFactory) ProviderName() string {
	return "INFOBASKET"
}

func (m FonbetSuperleagueStatsProviderFactory) Create() (ports.StatsProvider, error) {
	return stats_provider.NewInfobasketStatsProviderAdapter(
		infobasket.NewInfobasketClient("org"),
		m.ExternalId,
	), nil
}
