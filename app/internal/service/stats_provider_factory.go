package service

import (
	"IMP/app/internal/adapters/stats_provider_factory"
	"IMP/app/internal/ports"
	"errors"
	"strconv"
)

func NewStatsProvider(leagueName string, tournamentId uint, tournamentsRepo ports.TournamentsRepo) (ports.StatsProvider, error) {
	var factory ports.StatsProviderFactory

	switch leagueName {
	case "NBA":
		factory = stats_provider_factory.NbaStatsProviderFactory{}
	case "MLBL":
		externalIdModel, err := tournamentsRepo.FindTournamentExternalId(tournamentId, stats_provider_factory.MlblStatsProviderFactory{}.ProviderName())
		if err != nil {
			return nil, err
		}

		externalId, err := strconv.Atoi(externalIdModel.ExternalId)
		if err != nil {
			return nil, err
		}

		factory = stats_provider_factory.MlblStatsProviderFactory{
			ExternalId: externalId,
		}
	default:
		return nil, errors.New("unknown league: " + leagueName)
	}

	return factory.Create()
}
