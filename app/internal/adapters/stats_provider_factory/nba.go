package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/infra/api_nba"
	"IMP/app/internal/ports"
	"os"
)

type NbaStatsProviderFactory struct{}

func (m NbaStatsProviderFactory) ProviderName() string {
	return "API_SPORTS"
}

func (m NbaStatsProviderFactory) Create() (ports.StatsProvider, error) {
	apiNbaToken := os.Getenv("API_SPORT_API_KEY")

	return stats_provider.NewApiNbaStatsProviderAdapter(
		api_nba.NewClient("https://v2.nba.api-sports.io", apiNbaToken),
	), nil
}
