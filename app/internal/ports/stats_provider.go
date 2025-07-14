package ports

import "IMP/app/internal/core/games"

type StatsProviderPort interface {
	GetGamesStatsByDate() ([]games.GameStatEntity, error)
}
