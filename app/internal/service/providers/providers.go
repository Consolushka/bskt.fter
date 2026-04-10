package providers

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/infra/api_nba"
	"IMP/app/internal/infra/config"
	"IMP/app/internal/infra/infobasket"
	"IMP/app/internal/infra/sportoteka"
	"IMP/app/internal/ports"
	"errors"
	"fmt"
	"strconv"
)

type Provider string

const (
	ApiNba     = "API_NBA"
	CdnNba     = "CDN_NBA"
	Infobasket = "INFOBASKET"
	Sportoteka = "SPORTOTEKA"
)

func NewProvider(providerName string, externalId *string, params *map[string]interface{}, cfg config.ProvidersConfig) (ports.StatsProvider, error) {
	switch providerName {
	case ApiNba:
		return stats_provider.NewApiNbaStatsProviderAdapter(
			api_nba.NewClient("https://v2.nba.api-sports.io", cfg.ApiSportApiKey, cfg.ApiNbaRateLimit),
		), nil
	case CdnNba:
		return nil, errors.New("not implemented")
	case Infobasket:
		if externalId == nil {
			return nil, errors.New("external id must be set")
		}

		intExternalId, err := strconv.Atoi(*externalId)
		if err != nil {
			return nil, fmt.Errorf("atoi %v returned error: %w", *externalId, err)
		}

		if params == nil {
			return nil, errors.New("params must be set")
		}
		leadHost := (*params)["leadHost"]
		if leadHost == nil {
			return nil, errors.New("leadHost must be set")
		}

		return stats_provider.NewInfobasketStatsProviderAdapter(
			infobasket.NewInfobasketClient(leadHost.(string), cfg.InfobasketRateLimit),
			intExternalId,
		), nil
	case Sportoteka:
		if externalId == nil {
			return nil, errors.New("external id must be set")
		}

		if params == nil {
			return nil, errors.New("params must be set")
		}
		year := (*params)["year"]
		if year == nil {
			return nil, errors.New("year must be set")
		}

		return stats_provider.NewSportotekaStatsProvider(
			sportoteka.NewClient(cfg.SportotekaRateLimit),
			*externalId,
			int(year.(float64)),
		), nil
	}

	return nil, errors.New("unknown provider: " + providerName)
}
