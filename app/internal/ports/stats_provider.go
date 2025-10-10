package ports

import (
	"IMP/app/internal/core/games"
	"time"
)

type StatsProvider interface {
	GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error)
}
