package ports

import (
	"IMP/app/internal/core/tournament_poll_logs"
)

type PollLogRepo interface {
	Create(log tournament_poll_logs.TournamentPollLogModel) (tournament_poll_logs.TournamentPollLogModel, error)
	GetLatestSuccess(tournamentId uint) (tournament_poll_logs.TournamentPollLogModel, error)
}
