package ports

import "IMP/app/internal/core/tournaments"

type TournamentsRepo interface {
	ListActiveTournaments() ([]tournaments.TournamentModel, error)
	ListTournamentsByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error)
	FindTournamentExternalId(tournamentId uint, providerName string) (tournaments.TournamentExternalIdModel, error)
}
