package ports

import (
	"IMP/app/internal/core/poll_watermarks"
)

type PollLogRepo interface {
	Create(log poll_watermarks.TournamentPollLogModel) (poll_watermarks.TournamentPollLogModel, error)
	GetLatestSuccess(tournamentId uint) (poll_watermarks.TournamentPollLogModel, error)
}
