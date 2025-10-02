package service

import (
	"IMP/app/internal/adapters/stats_provider_factory"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"errors"
	"strconv"
)

func NewTournamentStatsProvider(tournament tournaments.TournamentModel) (ports.StatsProvider, error) {
	var factory ports.StatsProviderFactory

	switch tournament.League.Alias {
	case "NBA":
		factory = stats_provider_factory.NbaStatsProviderFactory{}
	case "MLBL":
		providerName := stats_provider_factory.MlblStatsProviderFactory{}.ProviderName()

		var externalIdModelForProvider tournaments.TournamentExternalIdModel
		for _, model := range tournament.ExternalIds {
			if model.ProviderName == providerName {
				externalIdModelForProvider = model
			}
		}

		if externalIdModelForProvider.ExternalId == "" {
			return nil, errors.New(providerName + " external id not found")
		}

		externalId, err := strconv.Atoi(externalIdModelForProvider.ExternalId)
		if err != nil {
			return nil, err
		}

		factory = stats_provider_factory.MlblStatsProviderFactory{
			ExternalId: externalId,
		}
	case "UBA":
		providerName := stats_provider_factory.UbaStatsProviderFactory{}.ProviderName()

		var externalIdModelForProvider tournaments.TournamentExternalIdModel
		for _, model := range tournament.ExternalIds {
			if model.ProviderName == providerName {
				externalIdModelForProvider = model
			}
		}

		if externalIdModelForProvider.ExternalId == "" {
			return nil, errors.New(providerName + " external id not found")
		}

		factory = stats_provider_factory.UbaStatsProviderFactory{
			Tag:  externalIdModelForProvider.ExternalId,
			Year: tournament.EndAt.Year(),
		}
	case "FONBETSL":
		providerName := stats_provider_factory.FonbetSuperleagueStatsProviderFactory{}.ProviderName()

		var externalIdModelForProvider tournaments.TournamentExternalIdModel
		for _, model := range tournament.ExternalIds {
			if model.ProviderName == providerName {
				externalIdModelForProvider = model
			}
		}

		if externalIdModelForProvider.ExternalId == "" {
			return nil, errors.New(providerName + " external id not found")
		}

		externalId, err := strconv.Atoi(externalIdModelForProvider.ExternalId)
		if err != nil {
			return nil, err
		}

		factory = stats_provider_factory.FonbetSuperleagueStatsProviderFactory{
			ExternalId: externalId,
		}
	default:
		return nil, errors.New("unknown league: " + tournament.League.Alias)
	}

	return factory.Create()
}
