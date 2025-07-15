package ports

import (
	"IMP/app/internal/core/games"
	"time"
)

type StatsProvider interface {
	GetGamesStatsByDate(date time.Time) ([]games.GameStatEntity, error)
}
