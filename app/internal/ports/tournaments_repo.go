package ports

import "IMP/app/internal/core/tournaments"

type TournamentsRepo interface {
	ListActive() ([]tournaments.TournamentModel, error)
	ListByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error)
	Get(id uint) (tournaments.TournamentModel, error)
	// Create атомарно сохраняет турнир вместе с его привязкой к провайдеру:
	// турнир без provider-маппинга бесполезен, поэтому пишем их в одной транзакции.
	Create(tournament tournaments.TournamentModel, provider tournaments.TournamentProvider) (tournaments.TournamentModel, error)
}
