package ports

import "IMP/app/internal/core/tournaments"

type TournamentsRepo interface {
	ListActive() ([]tournaments.TournamentModel, error)
	ListByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error)
	Get(id uint) (tournaments.TournamentModel, error)
}
