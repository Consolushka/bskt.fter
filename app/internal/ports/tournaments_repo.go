package ports

import "IMP/app/internal/core/tournaments"

type TournamentsRepo interface {
	ListActiveTournaments() ([]tournaments.TournamentModel, error)
}
