package ports

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"time"
)

type StatsProvider interface {
	GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error)
	GetPlayerBio(id string) (players.PlayerBioEntity, error)
}
